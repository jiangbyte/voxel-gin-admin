// internal/modules/iam/account/register_portal.go 门户注册（对齐 voxel-boot registerPortal）。
//
// Author: Charlie

package account

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/platform/idgen"
	"voxel-gin-admin/internal/domain/auth"
	"voxel-gin-admin/internal/domain/profile"
)

// IdentityExists 身份标识是否已被占用。
func (s *Service) IdentityExists(ctx context.Context, identityType, identifier string) bool {
	id := strings.TrimSpace(identifier)
	if id == "" {
		return false
	}
	if identityType == IdentityEmail {
		id = strings.ToLower(id)
	}
	_, err := s.repo.FindIdentity(ctx, identityType, id)
	return err == nil
}

// AllocateUniqueAccount 分配唯一登录账号名（对齐 voxel-boot allocateUniqueAccount）。
func (s *Service) AllocateUniqueAccount(ctx context.Context, base string) (string, error) {
	candidate := strings.TrimSpace(base)
	if candidate == "" {
		return "", fmt.Errorf("账号名无效")
	}
	suffix := 0
	for {
		if !s.IdentityExists(ctx, IdentityAccount, candidate) {
			return candidate, nil
		}
		suffix++
		candidate = base + strconv.Itoa(suffix)
		if len(candidate) > 64 {
			trim := base
			if len(trim) > 64-len(strconv.Itoa(suffix)) {
				trim = trim[:max(3, 64-len(strconv.Itoa(suffix)))]
			}
			candidate = trim + strconv.Itoa(suffix)
		}
	}
}

// RegisterPortal 创建门户账号、资料与默认角色/部门。
func (s *Service) RegisterPortal(ctx context.Context, in auth.PortalRegisterInput) (accountID, accountLogin string, err error) {
	accountLogin, err = security.RequireAccountLogin(in.AccountLogin)
	if err != nil {
		return "", "", err
	}
	if s.IdentityExists(ctx, IdentityAccount, accountLogin) {
		return "", "", fmt.Errorf("账号已存在")
	}
	accID := idgen.Next()
	acc := Account{
		ID: accID, PasswordHash: in.PasswordHash, AccountType: string(security.AccountPortal),
		AccountStatus: security.AccountStatusEnabled,
	}
	ident := Identity{
		ID: idgen.Next(), AccountID: accID, IdentityType: IdentityAccount, Identifier: accountLogin,
		Verified: 1, IsPrimary: 1, BindStatus: BindBound,
	}
	if err := s.repo.CreateAccount(ctx, acc, ident); err != nil {
		return "", "", err
	}
	if err := s.replaceAccountIdentities(ctx, accID,
		in.EmailEnabled, in.Email, in.Email, in.EmailVerified, BindBound,
		in.PhoneEnabled, in.Phone, in.Phone, in.PhoneVerified, BindBound); err != nil {
		return "", "", err
	}
	nick := in.Nickname
	if nick == nil || strings.TrimSpace(*nick) == "" {
		suffix := accID
		if len(accID) > 8 {
			suffix = accID[len(accID)-8:]
		}
		generated := "user-" + suffix
		nick = &generated
	}
	if err := s.portal.UpsertProfile(ctx, &profile.Profile{
		AccountID: accID, Nickname: nick, Phone: in.Phone, Email: in.Email,
	}); err != nil {
		return "", "", err
	}
	if err := s.AssignRegisterDefaults(ctx, accID, security.AccountPortal); err != nil {
		return "", "", err
	}
	return accID, accountLogin, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

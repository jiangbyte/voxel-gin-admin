// internal/modules/iam/account/scope.go 账号数据范围校验（对齐 voxel-boot DataScopeResolver.assertAccountAccessible）。
//
// Author: Charlie

package account

import (
	"context"

	contextx "voxel-gin-admin/internal/infrastructure/core/context"
	"voxel-gin-admin/internal/infrastructure/core/security/datascope"
)

const accountPagePerm = "iam:account:page"

func (s *Service) assertAccountAccessible(ctx context.Context, accountID string) error {
	sess := contextx.Session(ctx)
	if err := datascope.AssertAccountAccessibleMsg(ctx, s.repo.DB(), sess, accountID, accountPagePerm); err != nil {
		return err
	}
	return nil
}
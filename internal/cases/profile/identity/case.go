// Package identitycase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package identitycase

import (
	domainpkg "voxel-gin-admin/internal/domain/profile/identity"
)

// IdentityCase 用例入口：嵌入领域 Service。
type IdentityCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *IdentityCase {
	if svc == nil {
		return nil
	}
	return &IdentityCase{Service: svc}
}

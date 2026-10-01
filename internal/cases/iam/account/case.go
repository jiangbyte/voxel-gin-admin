// Package accountcase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package accountcase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/account"
)

// AccountCase 用例入口：嵌入领域 Service。
type AccountCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *AccountCase {
	if svc == nil {
		return nil
	}
	return &AccountCase{Service: svc}
}

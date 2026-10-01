// Package rolecase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package rolecase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/role"
)

// RoleCase 用例入口：嵌入领域 Service。
type RoleCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *RoleCase {
	if svc == nil {
		return nil
	}
	return &RoleCase{Service: svc}
}

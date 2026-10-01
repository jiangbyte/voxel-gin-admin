// Package deptcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package deptcase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/dept"
)

// DeptCase 用例入口：嵌入领域 Service。
type DeptCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *DeptCase {
	if svc == nil {
		return nil
	}
	return &DeptCase{Service: svc}
}

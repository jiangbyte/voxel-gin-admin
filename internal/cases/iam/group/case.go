// Package groupcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package groupcase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/group"
)

// GroupCase 用例入口：嵌入领域 Service。
type GroupCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *GroupCase {
	if svc == nil {
		return nil
	}
	return &GroupCase{Service: svc}
}

// Package workspacecase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package workspacecase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/workspace"
)

// WorkspaceCase 用例入口：嵌入领域 Service。
type WorkspaceCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *WorkspaceCase {
	if svc == nil {
		return nil
	}
	return &WorkspaceCase{Service: svc}
}

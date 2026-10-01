// Package codegencase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package codegencase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/codegen"
)

// CodegenCase 用例入口：嵌入领域 Service。
type CodegenCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *CodegenCase {
	if svc == nil {
		return nil
	}
	return &CodegenCase{Service: svc}
}

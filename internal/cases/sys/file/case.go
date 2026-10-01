// Package filecase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package filecase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/file"
)

// FileCase 用例入口：嵌入领域 Service。
type FileCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *FileCase {
	if svc == nil {
		return nil
	}
	return &FileCase{Service: svc}
}

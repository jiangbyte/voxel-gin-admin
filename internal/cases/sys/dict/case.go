// Package dictcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package dictcase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/dict"
)

// DictCase 用例入口：嵌入领域 Service。
type DictCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *DictCase {
	if svc == nil {
		return nil
	}
	return &DictCase{Service: svc}
}

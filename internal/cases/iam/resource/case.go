// Package resourcecase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package resourcecase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/resource"
)

// ResourceCase 用例入口：嵌入领域 Service。
type ResourceCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *ResourceCase {
	if svc == nil {
		return nil
	}
	return &ResourceCase{Service: svc}
}

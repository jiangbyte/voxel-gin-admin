// Package profilecase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package profilecase

import (
	domainpkg "voxel-gin-admin/internal/domain/profile"
)

// ProfileCase 用例入口：嵌入领域 Service。
type ProfileCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *ProfileCase {
	if svc == nil {
		return nil
	}
	return &ProfileCase{Service: svc}
}

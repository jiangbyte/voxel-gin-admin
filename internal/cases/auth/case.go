// Package authcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package authcase

import (
	domainpkg "voxel-gin-admin/internal/domain/auth"
)

// AuthCase 用例入口：嵌入领域 Service。
type AuthCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *AuthCase {
	if svc == nil {
		return nil
	}
	return &AuthCase{Service: svc}
}

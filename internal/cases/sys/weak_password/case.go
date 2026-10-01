// Package weak_passwordcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package weak_passwordcase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/weak_password"
)

// WeakPasswordCase 用例入口：嵌入领域 Service。
type WeakPasswordCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *WeakPasswordCase {
	if svc == nil {
		return nil
	}
	return &WeakPasswordCase{Service: svc}
}

// Package clientcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package clientcase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/client"
)

// ClientCase 用例入口：嵌入领域 Service。
type ClientCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *ClientCase {
	if svc == nil {
		return nil
	}
	return &ClientCase{Service: svc}
}

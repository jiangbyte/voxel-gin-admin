// Package positioncase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package positioncase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/position"
)

// PositionCase 用例入口：嵌入领域 Service。
type PositionCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *PositionCase {
	if svc == nil {
		return nil
	}
	return &PositionCase{Service: svc}
}

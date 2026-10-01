// Package relationcase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package relationcase

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/relation"
)

// RelationCase 用例入口：嵌入领域 Service。
type RelationCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *RelationCase {
	if svc == nil {
		return nil
	}
	return &RelationCase{Service: svc}
}

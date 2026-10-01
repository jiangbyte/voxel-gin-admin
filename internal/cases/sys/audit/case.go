// Package auditcase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package auditcase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/audit"
)

// AuditCase 用例入口：嵌入领域 Service。
type AuditCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *AuditCase {
	if svc == nil {
		return nil
	}
	return &AuditCase{Service: svc}
}

// Package jobcase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package jobcase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/job"
)

// JobCase 用例入口：嵌入领域 Service。
type JobCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *JobCase {
	if svc == nil {
		return nil
	}
	return &JobCase{Service: svc}
}

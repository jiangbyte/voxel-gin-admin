// Package noticecase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package noticecase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/notice"
)

// NoticeCase 用例入口：嵌入领域 Service。
type NoticeCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *NoticeCase {
	if svc == nil {
		return nil
	}
	return &NoticeCase{Service: svc}
}

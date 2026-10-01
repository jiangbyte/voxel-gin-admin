// Package feedbackcase 用例编排（具体类型，无接口）。
//
// Author: Charlie
package feedbackcase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/feedback"
)

// FeedbackCase 用例入口：嵌入领域 Service。
type FeedbackCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *FeedbackCase {
	if svc == nil {
		return nil
	}
	return &FeedbackCase{Service: svc}
}

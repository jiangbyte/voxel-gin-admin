// Package configcase 用例编排（xfg-ddd cases；具体类型，无接口）。
//
// Author: Charlie
package configcase

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/config"
)

// ConfigCase 用例入口：嵌入领域 Service。
type ConfigCase struct {
	*domainpkg.Service
}

// Wrap 包装领域 Service 为 Case。
func Wrap(svc *domainpkg.Service) *ConfigCase {
	if svc == nil {
		return nil
	}
	return &ConfigCase{Service: svc}
}

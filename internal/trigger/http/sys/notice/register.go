// Package noticehttp 模块自注册。
//
// Author: Charlie
package noticehttp

import (
	noticecase "voxel-gin-admin/internal/cases/sys/notice"
	domainpkg "voxel-gin-admin/internal/domain/sys/notice"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.notice", 60, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(noticecase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

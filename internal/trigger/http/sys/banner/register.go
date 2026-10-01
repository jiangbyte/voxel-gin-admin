// Package bannerhttp 模块自注册。
//
// Author: Charlie
package bannerhttp

import (
	bannercase "voxel-gin-admin/internal/cases/sys/banner"
	domainpkg "voxel-gin-admin/internal/domain/sys/banner"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.banner", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(bannercase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

// Package confighttp 模块自注册。
//
// Author: Charlie
package confighttp

import (
	configcase "voxel-gin-admin/internal/cases/sys/config"
	domainpkg "voxel-gin-admin/internal/domain/sys/config"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.config", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(configcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

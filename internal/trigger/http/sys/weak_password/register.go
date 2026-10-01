// Package weak_passwordhttp 模块自注册。
//
// Author: Charlie
package weak_passwordhttp

import (
	weak_passwordcase "voxel-gin-admin/internal/cases/sys/weak_password"
	domainpkg "voxel-gin-admin/internal/domain/sys/weak_password"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.weak_password", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(weak_passwordcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

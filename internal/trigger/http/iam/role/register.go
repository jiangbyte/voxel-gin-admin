// Package rolehttp 模块自注册。
//
// Author: Charlie
package rolehttp

import (
	rolecase "voxel-gin-admin/internal/cases/iam/role"
	domainpkg "voxel-gin-admin/internal/domain/iam/role"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.role", 40, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(rolecase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

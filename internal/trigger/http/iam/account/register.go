// Package accounthttp 模块自注册。
//
// Author: Charlie
package accounthttp

import (
	accountcase "voxel-gin-admin/internal/cases/iam/account"
	domainpkg "voxel-gin-admin/internal/domain/iam/account"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.account", 20, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(accountcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

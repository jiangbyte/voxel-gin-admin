// Package clienthttp 模块自注册。
//
// Author: Charlie
package clienthttp

import (
	clientcase "voxel-gin-admin/internal/cases/iam/client"
	domainpkg "voxel-gin-admin/internal/domain/iam/client"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.client", 40, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(clientcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

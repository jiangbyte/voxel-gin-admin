// Package resourcehttp 模块自注册。
//
// Author: Charlie
package resourcehttp

import (
	resourcecase "voxel-gin-admin/internal/cases/iam/resource"
	domainpkg "voxel-gin-admin/internal/domain/iam/resource"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.resource", 40, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(resourcecase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

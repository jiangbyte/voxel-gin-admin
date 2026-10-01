// Package positionhttp 模块自注册。
//
// Author: Charlie
package positionhttp

import (
	positioncase "voxel-gin-admin/internal/cases/iam/position"
	domainpkg "voxel-gin-admin/internal/domain/iam/position"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.position", 40, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(positioncase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

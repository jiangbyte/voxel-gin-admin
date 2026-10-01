// Package grouphttp 模块自注册。
//
// Author: Charlie
package grouphttp

import (
	groupcase "voxel-gin-admin/internal/cases/iam/group"
	domainpkg "voxel-gin-admin/internal/domain/iam/group"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.group", 40, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(groupcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

// Package depthttp 模块自注册。
//
// Author: Charlie
package depthttp

import (
	deptcase "voxel-gin-admin/internal/cases/iam/dept"
	domainpkg "voxel-gin-admin/internal/domain/iam/dept"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.dept", 40, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(deptcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

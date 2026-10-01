// Package codegenhttp 模块自注册。
//
// Author: Charlie
package codegenhttp

import (
	codegencase "voxel-gin-admin/internal/cases/sys/codegen"
	domainpkg "voxel-gin-admin/internal/domain/sys/codegen"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.codegen", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(codegencase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

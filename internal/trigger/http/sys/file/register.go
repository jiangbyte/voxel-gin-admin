// Package filehttp 模块自注册。
//
// Author: Charlie
package filehttp

import (
	filecase "voxel-gin-admin/internal/cases/sys/file"
	domainpkg "voxel-gin-admin/internal/domain/sys/file"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.file", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(filecase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

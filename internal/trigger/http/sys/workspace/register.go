// Package workspacehttp 模块自注册。
//
// Author: Charlie
package workspacehttp

import (
	workspacecase "voxel-gin-admin/internal/cases/sys/workspace"
	domainpkg "voxel-gin-admin/internal/domain/sys/workspace"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.workspace", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(workspacecase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

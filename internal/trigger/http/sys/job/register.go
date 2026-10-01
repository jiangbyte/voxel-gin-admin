// Package jobhttp 模块自注册。
//
// Author: Charlie
package jobhttp

import (
	jobcase "voxel-gin-admin/internal/cases/sys/job"
	domainpkg "voxel-gin-admin/internal/domain/sys/job"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.job", 30, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(jobcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

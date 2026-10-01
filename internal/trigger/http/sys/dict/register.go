// Package dicthttp 模块自注册。
//
// Author: Charlie
package dicthttp

import (
	dictcase "voxel-gin-admin/internal/cases/sys/dict"
	domainpkg "voxel-gin-admin/internal/domain/sys/dict"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.dict", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(dictcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

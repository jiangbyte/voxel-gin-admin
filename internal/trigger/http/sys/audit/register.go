// Package audithttp 模块自注册。
//
// Author: Charlie
package audithttp

import (
	auditcase "voxel-gin-admin/internal/cases/sys/audit"
	domainpkg "voxel-gin-admin/internal/domain/sys/audit"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.audit", 50, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(auditcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

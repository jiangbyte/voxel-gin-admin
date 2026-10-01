// Package identityhttp 模块自注册。
//
// Author: Charlie
package identityhttp

import (
	identitycase "voxel-gin-admin/internal/cases/profile/identity"
	domainpkg "voxel-gin-admin/internal/domain/profile/identity"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("profile.identity", 69, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(identitycase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

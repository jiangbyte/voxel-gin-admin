// Package authhttp 模块自注册：组装 Case，路由仍挂在同包领域 Service（auth HTTP 与 oauth/session 同包）。
//
// Author: Charlie
package authhttp

import (
	authcase "voxel-gin-admin/internal/cases/auth"
	domainpkg "voxel-gin-admin/internal/domain/auth"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

const accountFinderKey = "account_finder"

func init() {
	module.Register("auth", 25, func(d *module.Deps) module.Module {
		var finder domainpkg.AccountFinder
		if v, ok := d.Service(accountFinderKey); ok {
			finder = v.(domainpkg.AccountFinder)
		}
		m := domainpkg.New(d, finder)
		svc := domainpkg.ServiceFromDeps(d, finder)
		_ = authcase.Wrap(svc)
		m.Routes = []module.RouteRegistrar{svc.RegisterHTTPRoutes}
		return m
	})
}

// Package admin 管理端用户中心模块注册（对齐 voxel-boot/fastapi profile.admin）。
//
// Author: Charlie
package admin

import (
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/platform/module"
	"voxel-gin-admin/internal/domain/profile"
	"voxel-gin-admin/internal/domain/profile/identity"
	"voxel-gin-admin/internal/domain/sys/file"
)

// init 自注册 profile.admin。
func init() {
	module.Register("profile.admin", 70, func(d *module.Deps) module.Module {
		idSvc := identity.FromDeps(d)
		s := profile.NewService(d.DB, d.Redis, d.Notify, d.Storage, file.FromDeps(d), d.Runtime,
			d.Audit, security.AccountAdmin, profile.ProfileTableAdmin, "admin", idSvc)
		return module.Module{
			Name:   "profile.admin",
			Order:  70,
			Models: []any{&profile.AdminProfileModel{}},
			Routes: []module.RouteRegistrar{s.AdminRoutes},
		}
	})
}

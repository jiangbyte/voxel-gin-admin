// Package portal 门户用户中心模块注册（对齐 voxel-boot/fastapi profile.portal）。
//
// Author: Charlie
package portal

import (
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/platform/module"
	"voxel-gin-admin/internal/domain/profile"
	"voxel-gin-admin/internal/domain/profile/identity"
	"voxel-gin-admin/internal/domain/sys/file"
)

// init 自注册 profile.portal。
func init() {
	module.Register("profile.portal", 71, func(d *module.Deps) module.Module {
		idSvc := identity.FromDeps(d)
		_ = profile.NewService(d.DB, d.Redis, d.Notify, d.Storage, file.FromDeps(d), d.Runtime,
			d.Audit, security.AccountPortal, profile.ProfileTablePortal, "portal", idSvc)
		return module.Module{
			Name:   "profile.portal",
			Order:  71,
			Models: []any{&profile.PortalProfileModel{}},
			// 门户用户中心路由由 voxel-gin-portal 挂载，admin 进程不注册 /v1/portal/*
			Routes: nil,
		}
	})
}

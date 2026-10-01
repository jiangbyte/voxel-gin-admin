// Package identity 实名认证依赖装配（无 HTTP 注册；路由在 trigger）。
//
// Author: Charlie
package identity

import (
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// ServiceKey Deps 服务袋键：注册 *Service 单例供 profile 等消费。
const ServiceKey = "profile_identity_service"

// FromDeps 从依赖袋取出实名认证服务；缺失时按 Deps 新建并注册。
func FromDeps(d *module.Deps) *Service {
	if d == nil {
		return nil
	}
	if v, ok := d.Service(ServiceKey); ok {
		if s, ok := v.(*Service); ok && s != nil {
			return s
		}
	}
	s := mustNewService(d)
	d.Provide(ServiceKey, s)
	return s
}

// ServiceFromDeps 供 Case/Trigger 组装（与 FromDeps 同单例）。
func ServiceFromDeps(d *module.Deps) *Service {
	return FromDeps(d)
}

func mustNewService(d *module.Deps) *Service {
	fallbackKey := ""
	if d.Cfg != nil {
		fallbackKey = d.Cfg.Crypto.FernetKey
	}
	crypto, err := NewFieldCrypto(d.Runtime, fallbackKey)
	if err != nil {
		panic("profile.identity: " + err.Error())
	}
	return NewService(d.DB, crypto, d.Storage)
}

// New 仅贡献 Models（路由由 trigger 挂载）。
func New(d *module.Deps) module.Module {
	_ = FromDeps(d)
	return module.Module{
		Name:   "profile.identity",
		Order:  69,
		Models: []any{&ProfileIdentity{}, &RealNameCase{}, &RealNameCaseRecord{}},
	}
}

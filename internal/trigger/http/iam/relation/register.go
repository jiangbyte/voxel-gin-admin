// Package relationhttp 模块自注册（无 HTTP；仅领域 Models）。
//
// Author: Charlie
package relationhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/relation"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("iam.relation", 40, func(d *module.Deps) module.Module {
		return domainpkg.New(d)
	})
}

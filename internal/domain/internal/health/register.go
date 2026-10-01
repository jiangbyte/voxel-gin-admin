// internal/domain/internal/health/register.go 模块自注册。
//
// Author: Charlie

package health

import (
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// init 自注册 internal.health 模块。
func init() {
	module.Register("internal.health", 5, func(d *module.Deps) module.Module {
		return New(d)
	})
}

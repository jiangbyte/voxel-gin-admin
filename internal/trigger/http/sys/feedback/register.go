// Package feedbackhttp 模块自注册。
//
// Author: Charlie
package feedbackhttp

import (
	feedbackcase "voxel-gin-admin/internal/cases/sys/feedback"
	domainpkg "voxel-gin-admin/internal/domain/sys/feedback"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

func init() {
	module.Register("sys.feedback", 60, func(d *module.Deps) module.Module {
		m := domainpkg.New(d)
		svc := domainpkg.ServiceFromDeps(d)
		h := NewHandler(feedbackcase.Wrap(svc))
		m.Routes = []module.RouteRegistrar{h.registerRoutes(d)}
		return m
	})
}

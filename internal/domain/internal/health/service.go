// internal/modules/internal/health/service.go 业务服务。
//
// Author: Charlie

package health

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Service å¥åº·æ£€æŸ¥æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewService æž„é€ æœåŠ¡ã€‚
func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{db: db, redis: rdb}
}

// New æž„å»º internal.health æ¨¡å—ã€‚

// ServiceFromDeps 从模块依赖构造领域 Service（供 Case/Trigger 组装）。
func ServiceFromDeps(d *module.Deps) *Service {

	return NewService(d.DB, d.Redis)
}

func New(d *module.Deps) module.Module {
	s := ServiceFromDeps(d)
	return module.Module{
		Name:   "internal.health",
		Routes: []module.RouteRegistrar{s.registerRoutes},
	}
}

// Live å­˜æ´»æ£€æŸ¥ã€‚
func (s *Service) Live() LiveResult {
	return LiveResult{Status: "live"}
}

// Ready å°±ç»ªæ£€æŸ¥ã€‚
func (s *Service) Ready(ctx context.Context) (ReadyResult, bool) {
	out := ReadyResult{Status: "ready"}
	out.Checks.Database = CheckItem{Enabled: 1}
	out.Checks.Redis = CheckItem{Enabled: 1}

	if sqlDB, err := s.db.DB(); err != nil {
		detail := err.Error()
		out.Checks.Database.Detail = &detail
	} else if err := sqlDB.PingContext(ctx); err != nil {
		detail := err.Error()
		out.Checks.Database.Detail = &detail
	} else {
		out.Checks.Database.OK = 1
		d := "connection ok"
		out.Checks.Database.Detail = &d
	}

	if err := s.redis.Ping(ctx).Err(); err != nil {
		detail := err.Error()
		out.Checks.Redis.Detail = &detail
	} else {
		out.Checks.Redis.OK = 1
		d := "connection ok"
		out.Checks.Redis.Detail = &d
	}

	ok := out.Checks.Database.OK == 1 && out.Checks.Redis.OK == 1
	if !ok {
		out.Status = "not_ready"
	}
	return out, ok
}

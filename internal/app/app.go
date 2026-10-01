// Package app 是应用装配根：基础设施、自注册模块、HTTP 与内嵌任务调度器。
//
// 默认 blank import internal/modules/all；复杂场景可直接改 framework。
//
// Author: Charlie
package app

import (
	"strings"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"voxel-gin-admin/internal/infrastructure/core/config"
	"voxel-gin-admin/internal/infrastructure/core/crypto"
	"voxel-gin-admin/internal/infrastructure/core/logger"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/audit"
	"voxel-gin-admin/internal/infrastructure/platform/cache"
	"voxel-gin-admin/internal/infrastructure/platform/configsync"
	"voxel-gin-admin/internal/infrastructure/platform/db"
	"voxel-gin-admin/internal/infrastructure/platform/gojob"
	"voxel-gin-admin/internal/infrastructure/platform/idgen"
	"voxel-gin-admin/internal/infrastructure/platform/module"
	"voxel-gin-admin/internal/infrastructure/platform/notify"
	"voxel-gin-admin/internal/infrastructure/platform/otel"
	"voxel-gin-admin/internal/infrastructure/platform/runtimecfg"
	"voxel-gin-admin/internal/infrastructure/platform/storage"
)

// Deps 应用进程级依赖。
//
// Author: Charlie
type Deps struct {
	Cfg      *config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *security.SessionStore
	Perms    *security.PermissionRegistry
	Storage  *storage.Manager
	Audit    *audit.Queue
	Notify   *notify.Facade
	Runtime  *runtimecfg.Settings
	Modules  *module.Registry
	Jobs     *gojob.Manager
}

// API 单体进程：HTTP + 模块钩子 + 内嵌任务调度器。
//
// Author: Charlie
type API struct {
	Deps             *Deps
	Engine           *gin.Engine
	Server           *http.Server
	Audit            *audit.Queue
	Jobs             *gojob.Manager
	configSyncCancel context.CancelFunc
}

// OpenInfra 连接 DB/Redis/存储，准备空 Deps（随后 AttachRegisteredModules）。
func OpenInfra(cfg *config.Config) (*Deps, error) {
	if err := logger.Setup(cfg.App.Debug); err != nil {
		return nil, err
	}
	if err := idgen.Init(cfg.IDGen.WorkerID, cfg.IDGen.DatacenterID); err != nil {
		return nil, err
	}
	gdb, err := db.Open(cfg.DB)
	if err != nil {
		return nil, err
	}
	rdb, err := cache.Open(cfg.Redis)
	if err != nil {
		return nil, err
	}
	store := storage.NewManager()
	if err := otel.Init(cfg.OTel); err != nil {
		return nil, err
	}
	nf := notify.NewFacade(cfg.Notify, gdb)
	rt := runtimecfg.New(gdb)
	if codec, err := crypto.NewFernetFromConfig(cfg.Crypto.FernetKey, cfg.Crypto.VaultAddr); err == nil {
		rt.WithCodec(codec)
	}
	store.SetRuntime(rt)
	return &Deps{
		Cfg:      cfg,
		DB:       gdb,
		Redis:    rdb,
		Sessions: security.NewSessionStore(rdb),
		Perms:    security.NewPermissionRegistry(rdb),
		Storage:  store,
		Audit:  audit.NewQueue(gdb, rdb, cfg.Audit),
		Notify: nf,
		Runtime:  rt,
		// 任务调度器（handlers 在 NewAPI 装配完成后填充）
		Jobs: gojob.NewManager(gdb, rdb, gojob.Config{
			ScanIntervalMS:   cfg.Job.ScanIntervalMS,
			PoolSize:         cfg.Job.PoolSize,
			LogRetentionDays: cfg.Job.LogRetentionDays,
			LogBatchSize:     cfg.Job.LogBatchSize,
		}, nil),
	}, nil
}

// NewAPI 构建 Gin 引擎与 HTTP Server。
func NewAPI(d *Deps) *API {
	if d.Cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(VoxelSideFilter("admin"))
	r.Use(middleware.Recovery())
	r.Use(middleware.SecurityHeaders(d.Cfg.Security))
	r.Use(middleware.AccessLog())
	r.Use(middleware.AuthContext(d.Cfg.Auth, d.Sessions))
	r.Use(middleware.CSRFProtect(d.Cfg.Auth))
	r.Use(middleware.AuthWhitelist(d.Cfg.Auth.AuthWhitelist))
	r.Use(middleware.Trace())
	r.Use(middleware.CORS(d.Cfg.CORS))
	r.Use(middleware.ErrorHandler())

	r.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{"name": d.Cfg.App.Name})
	})
	MountOpenAPI(r, "http://"+d.Cfg.Addr(), d.Cfg.App.Name)
	if d.Cfg.Metrics.Enabled {
		path := d.Cfg.Metrics.Path
		if path == "" {
			path = "/metrics"
		}
		r.GET(path, middleware.PrometheusHandler())
	}

	api := r.Group("/api")
	// 操作审计由各路由挂载 middleware.OperationAudit(d.Audit, resourceType, action)
	d.Modules.MountRoutes(api)

	srv := &http.Server{
		Addr:              d.Cfg.Addr(),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	// 模块装配完成后填充任务处理器
	if d.Jobs != nil {
		d.Jobs.SetHandlers(collectHandlers(d.Modules))
	}
	return &API{
		Deps:   d,
		Engine: r,
		Server: srv,
		Audit:  d.Audit,
		Jobs:   d.Jobs,
	}
}

// collectHandlers 从模块注册表提取任务处理器（module.Job → gojob.HandlerDef）。
func collectHandlers(regs *module.Registry) []gojob.HandlerDef {
	if regs == nil {
		return nil
	}
	var handlers []gojob.HandlerDef
	for _, mod := range regs.Modules {
		for _, j := range mod.Jobs {
			j := j
			handlers = append(handlers, gojob.HandlerDef{Key: j.Name, Name: j.Name, Run: j.Run})
		}
	}
	return handlers
}

// Start 启动审计队列、模块钩子、任务调度器与 HTTP 监听。
func (a *API) Start(ctx context.Context) error {
	a.configSyncCancel = configsync.Start(ctx, a.Deps.Redis, func() {
		if a.Deps.Runtime != nil {
			a.Deps.Runtime.Invalidate()
		}
		if a.Deps.Storage != nil {
			a.Deps.Storage.Refresh()
		}
	})
	a.Audit.Start(ctx)
	if err := a.Deps.Modules.RunStart(ctx); err != nil {
		return err
	}
	if err := a.Deps.Perms.Sync(ctx); err != nil {
		logger.L.Warn("权限注册表同步失败", zap.Error(err))
	}
	if a.Jobs != nil {
		if err := a.Jobs.Start(ctx); err != nil {
			return err
		}
		logger.L.Info("任务调度器已启动")
	}
	logger.L.Info("api 正在启动", zap.String("addr", a.Server.Addr))
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.Server.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-time.After(200 * time.Millisecond):
		logger.L.Info("api 已监听", zap.String("addr", a.Server.Addr))
	}
	go func() {
		if err := <-errCh; err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("监听失败", zap.Error(err))
		}
	}()
	return nil
}

// Stop 优雅关闭。
func (a *API) Stop(ctx context.Context) error {
	if a.configSyncCancel != nil {
		a.configSyncCancel()
		a.configSyncCancel = nil
	}
	if a.Jobs != nil {
		_ = a.Jobs.Stop(ctx)
	}
	_ = a.Deps.Modules.RunStop(ctx)
	a.Audit.Stop()
	err := a.Server.Shutdown(ctx)
	_ = a.Deps.Redis.Close()
	sqlDB, e := a.Deps.DB.DB()
	if e == nil {
		_ = sqlDB.Close()
	}
	logger.Sync()
	return err
}

// CloseIdle 关闭空闲连接。
func CloseIdle(d *Deps) error {
	if d == nil {
		return nil
	}
	var err error
	if d.Redis != nil {
		err = d.Redis.Close()
	}
	if d.DB != nil {
		sqlDB, e := d.DB.DB()
		if e == nil {
			if e2 := sqlDB.Close(); e2 != nil && err == nil {
				err = e2
			}
		}
	}
	return err
}

// LoadOrDie 加载配置，失败则 panic。
func LoadOrDie(path string) *config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		panic(fmt.Errorf("config: %w", err))
	}
	return cfg
}


// VoxelSideFilter 按进程端拒绝对端路由（admin/portal 拆进程）。
func VoxelSideFilter(side string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if side == "admin" && strings.Contains(p, "/v1/portal/") {
			c.AbortWithStatusJSON(404, gin.H{"code": "404", "info": "not found"})
			return
		}
		if side == "portal" && strings.Contains(p, "/v1/admin/") {
			c.AbortWithStatusJSON(404, gin.H{"code": "404", "info": "not found"})
			return
		}
		c.Next()
	}
}

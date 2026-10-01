// Package jobhttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package jobhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/job"

	jobcase "voxel-gin-admin/internal/cases/sys/job"

	"net/http"

	"github.com/gin-gonic/gin"

	"voxel-gin-admin/internal/infrastructure/core/bind"
	contextx "voxel-gin-admin/internal/infrastructure/core/context"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/types/schema"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/gojob"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Handler HTTP 适配：只调用 Case。
type Handler struct {
	uc *jobcase.JobCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *jobcase.JobCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/job/page", admin, middleware.RequirePermission(d.Perms, "sys:job:page", "任务分页"), h.page)
		api.GET("/v1/admin/sys/job/detail", admin, middleware.RequirePermission(d.Perms, "sys:job:detail", "任务详情"), h.detail)
		api.POST("/v1/admin/sys/job/create", admin, middleware.RequirePermission(d.Perms, "sys:job:create", "任务创建"), middleware.OperationAudit(d.Audit, "sys_job", "create"), h.create)
		api.POST("/v1/admin/sys/job/update", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务更新"), middleware.OperationAudit(d.Audit, "sys_job", "update"), h.update)
		api.POST("/v1/admin/sys/job/delete", admin, middleware.RequirePermission(d.Perms, "sys:job:delete", "任务删除"), middleware.OperationAudit(d.Audit, "sys_job", "delete"), h.delete)
		api.POST("/v1/admin/sys/job/enabled", admin, middleware.RequirePermission(d.Perms, "sys:job:update", "任务启停"), middleware.OperationAudit(d.Audit, "sys_job", "enabled"), h.enabled)
		api.POST("/v1/admin/sys/job/run", admin, middleware.RequirePermission(d.Perms, "sys:job:run", "任务立即执行"), middleware.OperationAudit(d.Audit, "sys_job", "run"), h.run)
		api.GET("/v1/admin/sys/job-log/page", admin, middleware.RequirePermission(d.Perms, "sys:joblog:page", "任务日志分页"), h.logs)
	}
}

func (h *Handler) page(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.Page(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (h *Handler) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := h.uc.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) create(c *gin.Context) {
	var req domainpkg.AddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Create(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) update(c *gin.Context) {
	var req domainpkg.EditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Update(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) delete(c *gin.Context) {
	var body domainpkg.IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Delete(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) enabled(c *gin.Context) {
	var req domainpkg.EnabledParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.SetEnabled(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) run(c *gin.Context) {
	var req domainpkg.RunParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	executor := gojob.ExecutorSystem
	if sess := contextx.Session(c.Request.Context()); sess != nil && sess.AccountID != "" {
		executor = sess.AccountID
	}
	if err := h.uc.RunNow(c.Request.Context(), req.ID, executor); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) logs(c *gin.Context) {
	var q domainpkg.LogParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.Logs(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

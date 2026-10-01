// Package confighttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package confighttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/config"

	configcase "voxel-gin-admin/internal/cases/sys/config"

	"net/http"

	"github.com/gin-gonic/gin"

	"voxel-gin-admin/internal/infrastructure/core/bind"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/types/schema"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Handler HTTP 适配：只调用 Case。
type Handler struct {
	uc *configcase.ConfigCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *configcase.ConfigCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/config/create", admin, middleware.RequirePermission(d.Perms, "sys:config:create", "配置创建"), middleware.OperationAudit(d.Audit, "sys_config", "create"), h.create)
		api.POST("/v1/admin/sys/config/update", admin, middleware.RequirePermission(d.Perms, "sys:config:update", "配置更新"), middleware.OperationAudit(d.Audit, "sys_config", "update"), h.update)
		api.POST("/v1/admin/sys/config/delete", admin, middleware.RequirePermission(d.Perms, "sys:config:delete", "配置删除"), middleware.OperationAudit(d.Audit, "sys_config", "delete"), h.delete)
		api.GET("/v1/admin/sys/config/detail", admin, middleware.RequirePermission(d.Perms, "sys:config:detail", "配置详情"), h.detail)
		api.GET("/v1/admin/sys/config/page", admin, middleware.RequirePermission(d.Perms, "sys:config:page", "配置分页"), h.page)
		api.GET("/v1/admin/sys/config/list", admin, middleware.RequirePermission(d.Perms, "sys:config:page", "配置列表"), h.list)
		api.POST("/v1/admin/sys/config/batch-save", admin, middleware.RequirePermission(d.Perms, "sys:config:update", "配置批量保存"), middleware.OperationAudit(d.Audit, "sys_config", "update"), h.batchSave)
		api.POST("/v1/admin/sys/config/audit-alert/test-webhook", admin, middleware.OperationAudit(d.Audit, "sys_config", "test_webhook"), h.testWebhook)
		api.POST("/v1/admin/sys/config/audit-alert/test-push", admin, middleware.OperationAudit(d.Audit, "sys_config", "test_push"), h.testPush)
	}
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

func (h *Handler) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := h.uc.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
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

func (h *Handler) list(c *gin.Context) {
	var q domainpkg.ListParam
	q.Category = c.Query("category")
	q.Scope = c.Query("scope")
	rows, err := h.uc.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) batchSave(c *gin.Context) {
	var req domainpkg.BatchSaveParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.BatchSave(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) testWebhook(c *gin.Context) {
	var req domainpkg.TestWebhookParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	url := req.URL
	if url == "" {
		url = req.WebhookURL
	}
	if url == "" {
		response.Fail(c, http.StatusBadRequest, 400, "url required")
		return
	}
	secret := req.Secret
	if secret == "" {
		secret = req.WebhookSecret
	}
	if err := h.uc.TestWebhook(c.Request.Context(), url, secret); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, map[string]string{"message": "测试消息已发送"})
}

func (h *Handler) testPush(c *gin.Context) {
	if err := h.uc.TestPush(c.Request.Context()); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, map[string]string{"message": "测试消息已发送"})
}

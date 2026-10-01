// Package clienthttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package clienthttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/client"

	clientcase "voxel-gin-admin/internal/cases/iam/client"

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
	uc *clientcase.ClientCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *clientcase.ClientCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/client-module/create", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:create", "客户端模块创建"), middleware.OperationAudit(d.Audit, "iam_clientmodule", "create"), h.createModule)
		api.POST("/v1/admin/iam/client-module/update", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:update", "客户端模块更新"), middleware.OperationAudit(d.Audit, "iam_clientmodule", "update"), h.updateModule)
		api.POST("/v1/admin/iam/client-module/delete", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:delete", "客户端模块删除"), middleware.OperationAudit(d.Audit, "iam_clientmodule", "delete"), h.deleteModule)
		api.GET("/v1/admin/iam/client-module/detail", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:detail", "客户端模块详情"), h.detailModule)
		api.GET("/v1/admin/iam/client-module/page", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:page", "客户端模块分页"), h.pageModule)
		api.GET("/v1/admin/iam/client-module/selector", admin, middleware.RequirePermission(d.Perms, "iam:clientmodule:page", "客户端模块选择"), h.selectorModule)
		api.POST("/v1/admin/iam/client-resource/create", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:create", "客户端资源创建"), middleware.OperationAudit(d.Audit, "iam_clientresource", "create"), h.createResource)
		api.POST("/v1/admin/iam/client-resource/update", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:update", "客户端资源更新"), middleware.OperationAudit(d.Audit, "iam_clientresource", "update"), h.updateResource)
		api.POST("/v1/admin/iam/client-resource/delete", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:delete", "客户端资源删除"), middleware.OperationAudit(d.Audit, "iam_clientresource", "delete"), h.deleteResource)
		api.GET("/v1/admin/iam/client-resource/detail", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:detail", "客户端资源详情"), h.detailResource)
		api.GET("/v1/admin/iam/client-resource/page", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:page", "客户端资源分页"), h.pageResource)
		api.GET("/v1/admin/iam/client-resource/tree", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:list", "客户端资源树"), h.treeResource)
	}
}

func (h *Handler) createModule(c *gin.Context) {
	var req domainpkg.ModuleAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.CreateModule(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) updateModule(c *gin.Context) {
	var req domainpkg.ModuleEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.UpdateModule(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) deleteModule(c *gin.Context) {
	var body domainpkg.IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.DeleteModules(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) detailModule(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := h.uc.ModuleDetail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (h *Handler) pageModule(c *gin.Context) {
	var q domainpkg.ModulePageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.ModulePage(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (h *Handler) selectorModule(c *gin.Context) {
	out, err := h.uc.ModuleSelector(c.Request.Context(), c.Query("account_type"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}

func (h *Handler) createResource(c *gin.Context) {
	var req domainpkg.ResourceAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.CreateResource(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) updateResource(c *gin.Context) {
	var req domainpkg.ResourceEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.UpdateResource(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) deleteResource(c *gin.Context) {
	var body domainpkg.IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.DeleteResources(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) detailResource(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := h.uc.ResourceDetail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (h *Handler) pageResource(c *gin.Context) {
	var q domainpkg.ResourcePageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.ResourcePage(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (h *Handler) treeResource(c *gin.Context) {
	nodes, err := h.uc.ResourceTree(c.Request.Context(), c.Query("module_id"), c.Query("account_type"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}

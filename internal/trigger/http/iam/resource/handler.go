// Package resourcehttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package resourcehttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/resource"

	resourcecase "voxel-gin-admin/internal/cases/iam/resource"

	"net/http"

	"github.com/gin-gonic/gin"

	"voxel-gin-admin/internal/infrastructure/core/bind"
	contextx "voxel-gin-admin/internal/infrastructure/core/context"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/types/schema"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Handler HTTP 适配：只调用 Case。
type Handler struct {
	uc *resourcecase.ResourceCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *resourcecase.ResourceCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/resource/create", admin, middleware.RequirePermission(d.Perms, "iam:resource:create", "资源创建"), middleware.OperationAudit(d.Audit, "iam_resource", "create"), h.create)
		api.POST("/v1/admin/iam/resource/update", admin, middleware.RequirePermission(d.Perms, "iam:resource:update", "资源更新"), middleware.OperationAudit(d.Audit, "iam_resource", "update"), h.update)
		api.POST("/v1/admin/iam/resource/delete", admin, middleware.RequirePermission(d.Perms, "iam:resource:delete", "资源删除"), middleware.OperationAudit(d.Audit, "iam_resource", "delete"), h.delete)
		api.GET("/v1/admin/iam/resource/detail", admin, middleware.RequirePermission(d.Perms, "iam:resource:detail", "资源详情"), h.detail)
		api.GET("/v1/admin/iam/resource/page", admin, middleware.RequirePermission(d.Perms, "iam:resource:page", "资源分页"), h.page)
		api.GET("/v1/admin/iam/resource/current", admin, h.currentAdmin)
		api.GET("/v1/admin/iam/resource/tree", admin, middleware.RequirePermission(d.Perms, "iam:resource:list", "资源树"), h.tree)
		api.POST("/v1/admin/iam/resource-module/create", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:create", "资源模块创建"), middleware.OperationAudit(d.Audit, "iam_resourcemodule", "create"), h.createModule)
		api.POST("/v1/admin/iam/resource-module/update", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:update", "资源模块更新"), middleware.OperationAudit(d.Audit, "iam_resourcemodule", "update"), h.updateModule)
		api.POST("/v1/admin/iam/resource-module/delete", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:delete", "资源模块删除"), middleware.OperationAudit(d.Audit, "iam_resourcemodule", "delete"), h.deleteModule)
		api.GET("/v1/admin/iam/resource-module/detail", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:detail", "资源模块详情"), h.detailModule)
		api.GET("/v1/admin/iam/resource-module/page", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:page", "资源模块分页"), h.pageModule)
		api.GET("/v1/admin/iam/resource-module/selector", admin, middleware.RequirePermission(d.Perms, "iam:resourcemodule:page", "资源模块选择"), h.selectorModule)
		api.GET("/v1/admin/iam/permission/registry", admin, middleware.RequirePermission(d.Perms, "iam:resource:grant", "资源授权"), h.permissionRegistry)
		api.POST("/v1/admin/resource-permissions", admin, middleware.RequirePermission(d.Perms, "iam:resource:grant", "资源授权"), middleware.OperationAudit(d.Audit, "iam_resource", "grant"), h.bindResourcePermissions)
		api.POST("/v1/admin/client-resource-permissions", admin, middleware.RequirePermission(d.Perms, "iam:clientresource:grant", "客户端资源授权"), middleware.OperationAudit(d.Audit, "iam_clientresource", "grant"), h.bindClientResourcePermissions)
		api.POST("/v1/admin/iam/resource-button/create", admin, middleware.RequirePermission(d.Perms, "iam:resource:create", "资源创建"), middleware.OperationAudit(d.Audit, "iam_resource", "create"), h.createButton)
		api.POST("/v1/admin/iam/resource-button/update", admin, middleware.RequirePermission(d.Perms, "iam:resource:update", "资源更新"), middleware.OperationAudit(d.Audit, "iam_resource", "update"), h.updateButton)
		api.POST("/v1/admin/iam/resource-button/delete", admin, middleware.RequirePermission(d.Perms, "iam:resource:delete", "资源删除"), middleware.OperationAudit(d.Audit, "iam_resource", "delete"), h.deleteButton)
		api.GET("/v1/admin/iam/resource-button/page", admin, middleware.RequirePermission(d.Perms, "iam:resource:list", "资源分页"), h.pageButton)
	}
}

func (h *Handler) create(c *gin.Context) {
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

func (h *Handler) update(c *gin.Context) {
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

func (h *Handler) delete(c *gin.Context) {
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

func (h *Handler) detail(c *gin.Context) {
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

func (h *Handler) page(c *gin.Context) {
	var q domainpkg.ResourcePageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.ResourcePage(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (h *Handler) currentAdmin(c *gin.Context) {
	sess := contextx.Session(c.Request.Context())
	rows, err := h.uc.CurrentAdmin(c.Request.Context(), sessionAccountID(sess), sessionRoleIDs(sess), sessionGroupIDs(sess), sessionAllPerms(sess))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) currentPortal(c *gin.Context) {
	sess := contextx.Session(c.Request.Context())
	rows, err := h.uc.CurrentPortal(c.Request.Context(), sessionAccountID(sess), sessionRoleIDs(sess), sessionGroupIDs(sess), sessionAllPerms(sess))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func sessionAccountID(s *security.SessionPayload) string {
	if s == nil {
		return ""
	}
	return s.AccountID
}

func sessionRoleIDs(s *security.SessionPayload) []string {
	if s == nil {
		return nil
	}
	return s.RoleIDs
}

func sessionGroupIDs(s *security.SessionPayload) []string {
	if s == nil {
		return nil
	}
	return s.GroupIDs
}

func sessionAllPerms(s *security.SessionPayload) bool {
	if s == nil {
		return false
	}
	for _, k := range s.PermissionKeys {
		if k == "*:*:*" {
			return true
		}
	}
	return false
}

func (h *Handler) tree(c *gin.Context) {
	nodes, err := h.uc.ResourceTree(c.Request.Context(), c.Query("module_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
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
	out, err := h.uc.ModuleSelector(c.Request.Context(), c.Query("client"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}
func (h *Handler) permissionRegistry(c *gin.Context) {
	response.OK(c, h.uc.ListPermissions())
}

func (h *Handler) bindResourcePermissions(c *gin.Context) {
	var req domainpkg.ResourcePermissionBindParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.BindResourcePermissions(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) bindClientResourcePermissions(c *gin.Context) {
	var req domainpkg.ResourcePermissionBindParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.BindClientResourcePermissions(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) createButton(c *gin.Context) {
	var req domainpkg.ButtonAddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.CreateButton(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) updateButton(c *gin.Context) {
	var req domainpkg.ButtonEditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.UpdateButton(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) deleteButton(c *gin.Context) {
	var body domainpkg.IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.DeleteButtons(c.Request.Context(), body.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) pageButton(c *gin.Context) {
	var q domainpkg.ButtonPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.PageButtons(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

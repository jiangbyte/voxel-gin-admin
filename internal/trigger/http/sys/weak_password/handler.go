// Package weak_passwordhttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package weak_passwordhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/weak_password"

	weak_passwordcase "voxel-gin-admin/internal/cases/sys/weak_password"

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
	uc *weak_passwordcase.WeakPasswordCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *weak_passwordcase.WeakPasswordCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/weak-password/create", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:create", "弱密码创建"), middleware.OperationAudit(d.Audit, "sys_weakpassword", "create"), h.create)
		api.POST("/v1/admin/sys/weak-password/update", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:update", "弱密码更新"), middleware.OperationAudit(d.Audit, "sys_weakpassword", "update"), h.update)
		api.POST("/v1/admin/sys/weak-password/delete", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:delete", "弱密码删除"), middleware.OperationAudit(d.Audit, "sys_weakpassword", "delete"), h.delete)
		api.GET("/v1/admin/sys/weak-password/detail", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:detail", "弱密码详情"), h.detail)
		api.GET("/v1/admin/sys/weak-password/page", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:page", "弱密码分页"), h.page)
		api.GET("/v1/admin/sys/weak-password/list", admin, middleware.RequirePermission(d.Perms, "sys:weakpassword:list", "弱密码列表"), h.list)
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
	q.Password = c.Query("password")
	rows, err := h.uc.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

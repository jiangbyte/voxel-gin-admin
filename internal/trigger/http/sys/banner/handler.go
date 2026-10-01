// Package bannerhttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package bannerhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/banner"

	bannercase "voxel-gin-admin/internal/cases/sys/banner"

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
	uc *bannercase.BannerCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *bannercase.BannerCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/banner/list", admin, h.list)
		api.POST("/v1/admin/sys/banner/create", admin, middleware.RequirePermission(d.Perms, "sys:banner:create", "Banner创建"), middleware.OperationAudit(d.Audit, "sys_banner", "create"), h.create)
		api.POST("/v1/admin/sys/banner/update", admin, middleware.RequirePermission(d.Perms, "sys:banner:update", "Banner更新"), middleware.OperationAudit(d.Audit, "sys_banner", "update"), h.update)
		api.POST("/v1/admin/sys/banner/delete", admin, middleware.RequirePermission(d.Perms, "sys:banner:delete", "Banner删除"), middleware.OperationAudit(d.Audit, "sys_banner", "delete"), h.delete)
		api.GET("/v1/admin/sys/banner/detail", admin, middleware.RequirePermission(d.Perms, "sys:banner:detail", "Banner详情"), h.detail)
		api.GET("/v1/admin/sys/banner/page", admin, middleware.RequirePermission(d.Perms, "sys:banner:page", "Banner分页"), h.page)
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
	q.Position = c.Query("position")
	q.Category = c.Query("category")
	q.Type = c.Query("type")
	rows, err := h.uc.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) portalList(c *gin.Context) {
	var q domainpkg.PortalListParam
	q.Position = c.Query("position")
	q.Category = c.Query("category")
	q.Type = c.Query("type")
	rows, err := h.uc.PortalList(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) interaction(c *gin.Context) {
	var req domainpkg.InteractionParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Interaction(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, http.StatusNotFound, 404, "banner not found")
		return
	}
	response.OK(c, nil)
}

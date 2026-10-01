// Package dicthttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package dicthttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/dict"

	dictcase "voxel-gin-admin/internal/cases/sys/dict"

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
	uc *dictcase.DictCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *dictcase.DictCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/dict/create", admin, middleware.RequirePermission(d.Perms, "sys:dict:create", "字典创建"), middleware.OperationAudit(d.Audit, "sys_dict", "create"), h.create)
		api.POST("/v1/admin/sys/dict/update", admin, middleware.RequirePermission(d.Perms, "sys:dict:update", "字典更新"), middleware.OperationAudit(d.Audit, "sys_dict", "update"), h.update)
		api.POST("/v1/admin/sys/dict/delete", admin, middleware.RequirePermission(d.Perms, "sys:dict:delete", "字典删除"), middleware.OperationAudit(d.Audit, "sys_dict", "delete"), h.delete)
		api.GET("/v1/admin/sys/dict/detail", admin, middleware.RequirePermission(d.Perms, "sys:dict:detail", "字典详情"), h.detail)
		api.GET("/v1/admin/sys/dict/page", admin, middleware.RequirePermission(d.Perms, "sys:dict:page", "字典分页"), h.page)
		api.GET("/v1/admin/sys/dict/tree", admin, h.tree)
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

func (h *Handler) tree(c *gin.Context) {
	var q domainpkg.TreeParam
	q.Code = c.Query("code")
	q.Category = c.Query("category")
	nodes, err := h.uc.Tree(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}

func (h *Handler) portalTree(c *gin.Context) {
	var q domainpkg.TreeParam
	q.Category = c.Query("category")
	nodes, err := h.uc.Tree(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}

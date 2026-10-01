// Package depthttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package depthttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/dept"

	deptcase "voxel-gin-admin/internal/cases/iam/dept"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	uc *deptcase.DeptCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *deptcase.DeptCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/dept/create", admin, middleware.RequirePermission(d.Perms, "iam:dept:create", "部门创建"), middleware.OperationAudit(d.Audit, "iam_dept", "create"), h.create)
		api.POST("/v1/admin/iam/dept/update", admin, middleware.RequirePermission(d.Perms, "iam:dept:update", "部门更新"), middleware.OperationAudit(d.Audit, "iam_dept", "update"), h.update)
		api.POST("/v1/admin/iam/dept/delete", admin, middleware.RequirePermission(d.Perms, "iam:dept:delete", "部门删除"), middleware.OperationAudit(d.Audit, "iam_dept", "delete"), h.delete)
		api.GET("/v1/admin/iam/dept/detail", admin, middleware.RequirePermission(d.Perms, "iam:dept:detail", "部门详情"), h.detail)
		api.GET("/v1/admin/iam/dept/page", admin, middleware.RequirePermission(d.Perms, "iam:dept:page", "部门分页"), h.page)
		api.GET("/v1/admin/iam/dept/tree", admin, middleware.RequirePermission(d.Perms, "iam:dept:tree", "部门树"), h.tree)
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
	if err := h.uc.Update(c.Request.Context(), req, contextx.Session(c.Request.Context())); err != nil {
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
	if err := h.uc.Delete(c.Request.Context(), body.IDs, contextx.Session(c.Request.Context())); err != nil {
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
	row, err := h.uc.Detail(c.Request.Context(), q.ID, contextx.Session(c.Request.Context()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, 404, "not found")
			return
		}
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) page(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.Page(c.Request.Context(), q, contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (h *Handler) tree(c *gin.Context) {
	nodes, err := h.uc.Tree(c.Request.Context(), contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nodes)
}

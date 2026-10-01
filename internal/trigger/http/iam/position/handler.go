// Package positionhttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package positionhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/position"

	positioncase "voxel-gin-admin/internal/cases/iam/position"

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
	uc *positioncase.PositionCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *positioncase.PositionCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/position/create", admin, middleware.RequirePermission(d.Perms, "iam:position:create", "职位创建"), middleware.OperationAudit(d.Audit, "iam_position", "create"), h.create)
		api.POST("/v1/admin/iam/position/update", admin, middleware.RequirePermission(d.Perms, "iam:position:update", "职位更新"), middleware.OperationAudit(d.Audit, "iam_position", "update"), h.update)
		api.POST("/v1/admin/iam/position/delete", admin, middleware.RequirePermission(d.Perms, "iam:position:delete", "职位删除"), middleware.OperationAudit(d.Audit, "iam_position", "delete"), h.delete)
		api.GET("/v1/admin/iam/position/detail", admin, middleware.RequirePermission(d.Perms, "iam:position:detail", "职位详情"), h.detail)
		api.GET("/v1/admin/iam/position/page", admin, middleware.RequirePermission(d.Perms, "iam:position:page", "职位分页"), h.page)
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

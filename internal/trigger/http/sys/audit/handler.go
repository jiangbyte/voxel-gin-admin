// Package audithttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package audithttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/audit"

	auditcase "voxel-gin-admin/internal/cases/sys/audit"

	"net/http"

	"github.com/gin-gonic/gin"

	contextx "voxel-gin-admin/internal/infrastructure/core/context"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/types/schema"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Handler HTTP 适配：只调用 Case。
type Handler struct {
	uc *auditcase.AuditCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *auditcase.AuditCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET("/v1/admin/sys/audit/page", admin, middleware.RequirePermission(d.Perms, "sys:audit:page", "审计分页"), h.page)
		api.GET("/v1/admin/sys/audit/my-page", admin, h.myPage)
		api.GET("/v1/admin/sys/audit/detail", admin, middleware.RequirePermission(d.Perms, "sys:audit:detail", "审计详情"), h.detail)
		api.GET("/v1/admin/sys/audit/my-detail", admin, h.myDetail)
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
	response.Page(c, int64(cur), int64(size), total, domainpkg.ToOperationLogResults(rows))
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
	response.OK(c, domainpkg.ToOperationLogResult(*row))
}

func (h *Handler) myPage(c *gin.Context) {
	accountID := contextx.AccountID(c.Request.Context())
	if accountID == "" {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.MyPage(c.Request.Context(), accountID, q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, domainpkg.ToOperationLogResults(rows))
}

func (h *Handler) myDetail(c *gin.Context) {
	accountID := contextx.AccountID(c.Request.Context())
	if accountID == "" {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := h.uc.MyDetail(c.Request.Context(), accountID, q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, domainpkg.ToOperationLogResult(*row))
}

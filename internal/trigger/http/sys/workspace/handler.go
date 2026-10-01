// Package workspacehttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package workspacehttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/workspace"

	workspacecase "voxel-gin-admin/internal/cases/sys/workspace"

	"net/http"

	"github.com/gin-gonic/gin"

	"voxel-gin-admin/internal/infrastructure/core/bind"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Handler HTTP 适配：只调用 Case。
type Handler struct {
	uc *workspacecase.WorkspaceCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *workspacecase.WorkspaceCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.GET(
			"/v1/admin/sys/workspace/overview",
			admin,
			h.overview,
		)
		api.GET("/v1/admin/sys/workspace/shortcuts", admin, h.listShortcuts)
		api.POST(
			"/v1/admin/sys/workspace/shortcuts",
			admin,
			middleware.OperationAudit(d.Audit, "workspace_shortcut", "update"),
			h.replaceShortcuts,
		)
	}
}

func (h *Handler) overview(c *gin.Context) {
	out, err := h.uc.Overview(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (h *Handler) listShortcuts(c *gin.Context) {
	out, err := h.uc.ListShortcuts(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (h *Handler) replaceShortcuts(c *gin.Context) {
	var req domainpkg.ShortcutSaveParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := h.uc.ReplaceShortcuts(c.Request.Context(), req.ResourceIDs)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}

// Package rolehttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package rolehttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/role"

	rolecase "voxel-gin-admin/internal/cases/iam/role"

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
	uc *rolecase.RoleCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *rolecase.RoleCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/role/create", admin, middleware.RequirePermission(d.Perms, "iam:role:create", "角色创建"), middleware.OperationAudit(d.Audit, "iam_role", "create"), h.create)
		api.POST("/v1/admin/iam/role/update", admin, middleware.RequirePermission(d.Perms, "iam:role:update", "角色更新"), middleware.OperationAudit(d.Audit, "iam_role", "update"), h.update)
		api.POST("/v1/admin/iam/role/delete", admin, middleware.RequirePermission(d.Perms, "iam:role:delete", "角色删除"), middleware.OperationAudit(d.Audit, "iam_role", "delete"), h.delete)
		api.GET("/v1/admin/iam/role/detail", admin, middleware.RequirePermission(d.Perms, "iam:role:detail", "角色详情"), h.detail)
		api.GET("/v1/admin/iam/role/page", admin, middleware.RequirePermission(d.Perms, "iam:role:page", "角色分页"), h.page)
		api.GET("/v1/admin/iam/role/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:ownresource", "角色已拥有资源"), h.ownResource)
		api.POST("/v1/admin/iam/role/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:grantresource", "角色资源授权"), middleware.OperationAudit(d.Audit, "iam_role", "grant_resource"), h.grantResource)
		api.GET("/v1/admin/iam/role/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:ownclientresource", "角色已拥有客户端资源"), h.ownClientResource)
		api.POST("/v1/admin/iam/role/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:role:grantclientresource", "角色客户端资源授权"), middleware.OperationAudit(d.Audit, "iam_role", "grant_client_resource"), h.grantClientResource)
		api.GET("/v1/admin/iam/role/own-user", admin, middleware.RequirePermission(d.Perms, "iam:role:ownuser", "角色成员查询"), h.ownUser)
		api.POST("/v1/admin/iam/role/grant-user", admin, middleware.RequirePermission(d.Perms, "iam:role:grantuser", "角色成员授权"), middleware.OperationAudit(d.Audit, "iam_role", "grant_user"), h.grantUser)
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
func (h *Handler) ownResource(c *gin.Context) {
	var q domainpkg.OwnResourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnResources(c.Request.Context(), q.ID, q.AccountType)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (h *Handler) grantResource(c *gin.Context) {
	var req domainpkg.GrantResourceParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.GrantResources(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ownClientResource(c *gin.Context) {
	var q domainpkg.OwnResourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnClientResources(c.Request.Context(), q.ID, q.AccountType)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (h *Handler) grantClientResource(c *gin.Context) {
	var req domainpkg.GrantResourceParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.GrantClientResources(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ownUser(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnUsers(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (h *Handler) grantUser(c *gin.Context) {
	var req domainpkg.GrantUserParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.GrantUsers(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

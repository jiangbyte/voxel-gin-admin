// Package grouphttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package grouphttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/group"

	groupcase "voxel-gin-admin/internal/cases/iam/group"

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
	uc *groupcase.GroupCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *groupcase.GroupCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/group/create", admin, middleware.RequirePermission(d.Perms, "iam:group:create", "用户组创建"), middleware.OperationAudit(d.Audit, "iam_group", "create"), h.create)
		api.POST("/v1/admin/iam/group/update", admin, middleware.RequirePermission(d.Perms, "iam:group:update", "用户组更新"), middleware.OperationAudit(d.Audit, "iam_group", "update"), h.update)
		api.POST("/v1/admin/iam/group/delete", admin, middleware.RequirePermission(d.Perms, "iam:group:delete", "用户组删除"), middleware.OperationAudit(d.Audit, "iam_group", "delete"), h.delete)
		api.GET("/v1/admin/iam/group/detail", admin, middleware.RequirePermission(d.Perms, "iam:group:detail", "用户组详情"), h.detail)
		api.GET("/v1/admin/iam/group/page", admin, middleware.RequirePermission(d.Perms, "iam:group:page", "用户组分页"), h.page)
		api.GET("/v1/admin/iam/group/own-user", admin, middleware.RequirePermission(d.Perms, "iam:group:ownuser", "用户组成员查询"), h.ownUser)
		api.POST("/v1/admin/iam/group/grant-user", admin, middleware.RequirePermission(d.Perms, "iam:group:grantuser", "用户组成员授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_user"), h.grantUser)
		api.GET("/v1/admin/iam/group/own-role", admin, middleware.RequirePermission(d.Perms, "iam:group:ownrole", "用户组已拥有角色"), h.ownRole)
		api.POST("/v1/admin/iam/group/grant-role", admin, middleware.RequirePermission(d.Perms, "iam:group:grantrole", "用户组角色授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_role"), h.grantRole)
		api.GET("/v1/admin/iam/group/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:ownresource", "用户组已拥有资源"), h.ownResource)
		api.POST("/v1/admin/iam/group/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:grantresource", "用户组资源授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_resource"), h.grantResource)
		api.GET("/v1/admin/iam/group/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:ownclientresource", "用户组已拥有客户端资源"), h.ownClientResource)
		api.POST("/v1/admin/iam/group/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:group:grantclientresource", "用户组客户端资源授权"), middleware.OperationAudit(d.Audit, "iam_group", "grant_client_resource"), h.grantClientResource)
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

func (h *Handler) ownRole(c *gin.Context) {
	var q domainpkg.OwnResourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnRoles(c.Request.Context(), q.ID, q.AccountType)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (h *Handler) grantRole(c *gin.Context) {
	var req domainpkg.GrantRoleParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.GrantRoles(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
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

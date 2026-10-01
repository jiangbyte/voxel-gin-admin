// Package accounthttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package accounthttp

import (
	domainpkg "voxel-gin-admin/internal/domain/iam/account"

	accountcase "voxel-gin-admin/internal/cases/iam/account"

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
	uc *accountcase.AccountCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *accountcase.AccountCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/iam/account/create", admin, middleware.RequirePermission(d.Perms, "iam:account:create", "账户创建"), middleware.OperationAudit(d.Audit, "iam_account", "create"), h.create)
		api.POST("/v1/admin/iam/account/update", admin, middleware.RequirePermission(d.Perms, "iam:account:update", "账户更新"), middleware.OperationAudit(d.Audit, "iam_account", "update"), h.update)
		api.POST("/v1/admin/iam/account/update-login-identity", admin, middleware.RequirePermission(d.Perms, "iam:account:update", "账户登录身份更新"), middleware.OperationAudit(d.Audit, "iam_account", "update_login_identity"), h.updateLoginIdentity)
		api.POST("/v1/admin/iam/account/delete", admin, middleware.RequirePermission(d.Perms, "iam:account:delete", "账户删除"), middleware.OperationAudit(d.Audit, "iam_account", "delete"), h.delete)
		api.GET("/v1/admin/iam/account/detail", admin, middleware.RequirePermission(d.Perms, "iam:account:detail", "账户详情"), h.detail)
		api.GET("/v1/admin/iam/account/page", admin, middleware.RequirePermission(d.Perms, "iam:account:page", "账户分页"), h.page)
		api.GET("/v1/admin/iam/account/own-role", admin, middleware.RequirePermission(d.Perms, "iam:account:ownrole", "账号已拥有角色"), h.ownRole)
		api.POST("/v1/admin/iam/account/grant-role", admin, middleware.RequirePermission(d.Perms, "iam:account:grantrole", "账号角色授权"), middleware.OperationAudit(d.Audit, "iam_account", "grant_role"), h.grantRole)
		api.GET("/v1/admin/iam/account/own-group", admin, middleware.RequirePermission(d.Perms, "iam:account:owngroup", "账号已拥有用户组"), h.ownGroup)
		api.POST("/v1/admin/iam/account/grant-group", admin, middleware.RequirePermission(d.Perms, "iam:account:grantgroup", "账号用户组授权"), middleware.OperationAudit(d.Audit, "iam_account", "grant_group"), h.grantGroup)
		api.GET("/v1/admin/iam/account/own-dept", admin, middleware.RequirePermission(d.Perms, "iam:account:owndept", "账号已拥有部门"), h.ownDept)
		api.POST("/v1/admin/iam/account/grant-dept", admin, middleware.RequirePermission(d.Perms, "iam:account:grantdept", "账号部门授权"), middleware.OperationAudit(d.Audit, "iam_account", "grant_dept"), h.grantDept)
		api.GET("/v1/admin/iam/account/own-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:ownresource", "账号已拥有资源"), h.ownResource)
		api.POST("/v1/admin/iam/account/grant-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:grantresource", "账号资源授权"), middleware.OperationAudit(d.Audit, "iam_account", "grant_resource"), h.grantResource)
		api.GET("/v1/admin/iam/account/own-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:ownclientresource", "账号已拥有客户端资源"), h.ownClientResource)
		api.POST("/v1/admin/iam/account/grant-client-resource", admin, middleware.RequirePermission(d.Perms, "iam:account:grantclientresource", "账号客户端资源授权"), middleware.OperationAudit(d.Audit, "iam_account", "grant_client_resource"), h.grantClientResource)
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

func (h *Handler) updateLoginIdentity(c *gin.Context) {
	var req domainpkg.UpdateLoginIdentityParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.UpdateLoginIdentity(c.Request.Context(), req); err != nil {
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
	vo, err := h.uc.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, vo)
}

func (h *Handler) page(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	records, total, cur, size, err := h.uc.Page(c.Request.Context(), q, contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, records)
}
func (h *Handler) ownRole(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnRoles(c.Request.Context(), q.ID)
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

func (h *Handler) ownGroup(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnGroups(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (h *Handler) grantGroup(c *gin.Context) {
	var req domainpkg.GrantGroupParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.GrantGroups(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ownDept(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	vo, err := h.uc.OwnDepts(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, vo)
}

func (h *Handler) grantDept(c *gin.Context) {
	var req domainpkg.GrantDeptParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.GrantDepts(c.Request.Context(), req); err != nil {
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

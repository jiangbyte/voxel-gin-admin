// Package identityhttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie
package identityhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/profile/identity"

	identitycase "voxel-gin-admin/internal/cases/profile/identity"

	"errors"
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
	uc *identitycase.IdentityCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *identitycase.IdentityCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		adminUser := api.Group("/v1/admin", middleware.RequireAccountType(security.AccountAdmin))
		adminUser.GET("/profile/identity/status", h.adminIdentityStatus)
		adminUser.GET("/profile/identity/case/options", h.adminOptions)
		adminUser.POST("/profile/identity/case/submit", middleware.OperationAudit(d.Audit, "real_name_case", "submit"), h.adminSubmit)
		adminUser.POST("/profile/identity/case/init-third-party", middleware.OperationAudit(d.Audit, "real_name_case", "init_third_party"), h.adminInitThirdParty)
		adminUser.POST("/profile/identity/case/callback", h.adminCallback)
		adminUser.GET("/profile/identity/case/my-page", h.adminMyPage)

		manageCase := api.Group("/v1/admin/identity/case", middleware.RequireAccountType(security.AccountAdmin))
		manageCase.GET("/review-page", middleware.RequirePermission(d.Perms, "identity:case:verify", "实名审核分页"), h.reviewPage)
		manageCase.GET("/detail", middleware.RequirePermission(d.Perms, "identity:case:verify", "实名工单详情"), h.detail)
		manageCase.POST("/approve", middleware.RequirePermission(d.Perms, "identity:case:verify", "实名审核通过"), middleware.OperationAudit(d.Audit, "real_name_case", "approve"), h.approve)
		manageCase.POST("/reject", middleware.RequirePermission(d.Perms, "identity:case:verify", "实名审核驳回"), middleware.OperationAudit(d.Audit, "real_name_case", "reject"), h.reject)
		manageCred := api.Group("/v1/admin/identity/credential", middleware.RequireAccountType(security.AccountAdmin))
		manageCred.GET("/page", middleware.RequirePermission(d.Perms, "identity:credential:revoke", "实名快照分页"), h.identityPage)
		manageCred.POST("/revoke", middleware.RequirePermission(d.Perms, "identity:credential:revoke", "撤销实名"), middleware.OperationAudit(d.Audit, "profile_identity", "revoke"), h.revoke)
	}
}

func (h *Handler) adminIdentityStatus(c *gin.Context)   { h.identityStatus(c) }
func (h *Handler) portalIdentityStatus(c *gin.Context)  { h.identityStatus(c) }
func (h *Handler) adminOptions(c *gin.Context)          { h.options(c) }
func (h *Handler) portalOptions(c *gin.Context)         { h.options(c) }
func (h *Handler) adminSubmit(c *gin.Context)           { h.submit(c) }
func (h *Handler) portalSubmit(c *gin.Context)          { h.submit(c) }
func (h *Handler) adminInitThirdParty(c *gin.Context)   { h.initThirdParty(c) }
func (h *Handler) portalInitThirdParty(c *gin.Context)  { h.initThirdParty(c) }
func (h *Handler) adminCallback(c *gin.Context)         { h.callback(c) }
func (h *Handler) portalCallback(c *gin.Context)        { h.callback(c) }
func (h *Handler) adminMyPage(c *gin.Context)           { h.myPage(c) }
func (h *Handler) portalMyPage(c *gin.Context)          { h.myPage(c) }

func (h *Handler) identityStatus(c *gin.Context) {
	accountID := contextx.AccountID(c.Request.Context())
	out, err := h.uc.GetUserStatus(c.Request.Context(), accountID)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, out)
}

func (h *Handler) options(c *gin.Context) {
	response.OK(c, h.uc.Options(c.Request.Context()))
}

func (h *Handler) submit(c *gin.Context) {
	var req domainpkg.RealNameCaseSubmitParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	accountID := contextx.AccountID(c.Request.Context())
	if err := h.uc.Submit(c.Request.Context(), accountID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) initThirdParty(c *gin.Context) {
	var req domainpkg.RealNameCaseInitThirdPartyParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	accountID := contextx.AccountID(c.Request.Context())
	out, err := h.uc.InitThirdParty(c.Request.Context(), accountID, req)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, out)
}

func (h *Handler) callback(c *gin.Context) {
	var req domainpkg.RealNameCaseCallbackParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Callback(c.Request.Context(), req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) myPage(c *gin.Context) {
	var q domainpkg.RealNameCaseMyPageParam
	_ = c.ShouldBindQuery(&q)
	accountID := contextx.AccountID(c.Request.Context())
	rows, total, current, size, err := h.uc.MyPage(c.Request.Context(), accountID, q)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) reviewPage(c *gin.Context) {
	var q domainpkg.RealNameCaseReviewPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := h.uc.ReviewPage(c.Request.Context(), q)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := h.uc.Detail(c.Request.Context(), q.ID)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, out)
}

func (h *Handler) approve(c *gin.Context) {
	var req domainpkg.RealNameCaseApproveParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	reviewerID := contextx.AccountID(c.Request.Context())
	if err := h.uc.Approve(c.Request.Context(), reviewerID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) reject(c *gin.Context) {
	var req domainpkg.RealNameCaseRejectParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	reviewerID := contextx.AccountID(c.Request.Context())
	if err := h.uc.Reject(c.Request.Context(), reviewerID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) identityPage(c *gin.Context) {
	var q domainpkg.IdentityPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := h.uc.IdentityPage(c.Request.Context(), q)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) revoke(c *gin.Context) {
	var req domainpkg.IdentityRevokeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	operatorID := contextx.AccountID(c.Request.Context())
	if err := h.uc.Revoke(c.Request.Context(), operatorID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func failBiz(c *gin.Context, err error) {
	var be *domainpkg.BizError
	if errors.As(err, &be) {
		status := be.HTTPStatus
		if status == 0 {
			status = http.StatusBadRequest
		}
		code := be.Code
		if code == 0 {
			code = 400
		}
		response.Fail(c, status, code, be.Message)
		return
	}
	response.Fail(c, http.StatusBadRequest, 400, err.Error())
}

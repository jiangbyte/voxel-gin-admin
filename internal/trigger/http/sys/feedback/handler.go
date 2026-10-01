// Package feedbackhttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package feedbackhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/feedback"

	feedbackcase "voxel-gin-admin/internal/cases/sys/feedback"

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
	uc *feedbackcase.FeedbackCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *feedbackcase.FeedbackCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := api.Group("/v1/admin/sys/feedback", middleware.RequireAccountType(security.AccountAdmin))
		admin.GET("/page", middleware.RequirePermission(d.Perms, "sys:feedback:page", "domainpkg.Feedback page"), h.pageAdmin)
		admin.GET("/detail", middleware.RequirePermission(d.Perms, "sys:feedback:detail", "domainpkg.Feedback detail"), h.detail)
		admin.POST("/update", middleware.RequirePermission(d.Perms, "sys:feedback:update", "Reply feedback"), middleware.OperationAudit(d.Audit, "sys_feedback", "update"), h.update)
		admin.POST("/delete", middleware.RequirePermission(d.Perms, "sys:feedback:delete", "Delete feedback"), middleware.OperationAudit(d.Audit, "sys_feedback", "delete"), h.delete)
		admin.POST("/submit", middleware.OperationAudit(d.Audit, "sys_feedback", "submit"), h.submit)
		admin.GET("/my-page", h.myPage)
		admin.GET("/my-detail", h.myDetail)
	}
}

func (h *Handler) submit(c *gin.Context) {
	var req domainpkg.CreateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	meta := domainpkg.SubmitMeta{
		AccountType: string(sess.AccountType), AccountID: sess.AccountID, CreatedBy: sess.AccountID,
	}
	if err := h.uc.Submit(c.Request.Context(), req, meta); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) update(c *gin.Context) {
	var req domainpkg.UpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	aid := contextx.AccountID(c.Request.Context())
	if err := h.uc.Update(c.Request.Context(), req, domainpkg.ReplyMeta{RepliedBy: aid, UpdatedBy: aid}); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) delete(c *gin.Context) {
	var req domainpkg.IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Delete(c.Request.Context(), req.IDs); err != nil {
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
		response.Fail(c, http.StatusNotFound, 404, "feedback not found")
		return
	}
	response.OK(c, row)
}

func (h *Handler) pageAdmin(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := h.uc.PageAdmin(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) myPage(c *gin.Context) {
	var q schema.PageQuery
	_ = c.ShouldBindQuery(&q)
	sess := contextx.Session(c.Request.Context())
	rows, total, current, size, err := h.uc.MyPage(c.Request.Context(), q, sess.AccountID, string(sess.AccountType))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) myDetail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sess := contextx.Session(c.Request.Context())
	row, err := h.uc.MyDetail(c.Request.Context(), q.ID, sess.AccountID, string(sess.AccountType))
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "feedback not found")
		return
	}
	response.OK(c, row)
}

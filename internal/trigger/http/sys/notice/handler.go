// Package noticehttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package noticehttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/notice"

	noticecase "voxel-gin-admin/internal/cases/sys/notice"

	"net/http"
	"strings"
	"time"

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
	uc *noticecase.NoticeCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *noticecase.NoticeCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := api.Group("/v1/admin/sys/notice", middleware.RequireAccountType(security.AccountAdmin))
		admin.POST("/create", middleware.RequirePermission(d.Perms, "sys:notice:create", "创建 notice"), middleware.OperationAudit(d.Audit, "sys_notice", "create"), h.create)
		admin.POST("/update", middleware.RequirePermission(d.Perms, "sys:notice:update", "Update notice"), middleware.OperationAudit(d.Audit, "sys_notice", "update"), h.update)
		admin.POST("/delete", middleware.RequirePermission(d.Perms, "sys:notice:delete", "Delete notice"), middleware.OperationAudit(d.Audit, "sys_notice", "delete"), h.delete)
		admin.GET("/detail", middleware.RequirePermission(d.Perms, "sys:notice:detail", "domainpkg.Notice detail"), h.detail)
		admin.GET("/page", middleware.RequirePermission(d.Perms, "sys:notice:page", "domainpkg.Notice page"), h.pageAdmin)
		admin.POST("/publish", middleware.RequirePermission(d.Perms, "sys:notice:publish", "Publish notice"), middleware.OperationAudit(d.Audit, "sys_notice", "publish"), h.publish)
		admin.POST("/revoke", middleware.RequirePermission(d.Perms, "sys:notice:revoke", "Revoke notice"), middleware.OperationAudit(d.Audit, "sys_notice", "revoke"), h.revoke)
		admin.POST("/pin", middleware.RequirePermission(d.Perms, "sys:notice:pin", "Pin notice"), middleware.OperationAudit(d.Audit, "sys_notice", "pin"), h.pin)
		h.registerMyRoutes(admin)
	}
}

func (h *Handler) registerMyRoutes(g *gin.RouterGroup) {
	g.GET("/my-page", h.myPage)
	g.GET("/my-detail", h.myDetail)
	g.GET("/unread-count", h.unreadCount)
	g.POST("/read", h.markRead)
	g.POST("/read-all", h.markAllRead)
}

func (h *Handler) create(c *gin.Context) {
	var req domainpkg.CreateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	aid := contextx.AccountID(c.Request.Context())
	if err := h.uc.Create(c.Request.Context(), req, &aid, &aid); err != nil {
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
	if err := h.uc.Update(c.Request.Context(), req, &aid); err != nil {
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
		response.Fail(c, http.StatusNotFound, 404, "notice not found")
		return
	}
	response.OK(c, row)
}

func (h *Handler) pageAdmin(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := h.uc.PageAdmin(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) publish(c *gin.Context) {
	var req domainpkg.IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	now := time.Now().UTC()
	aid := contextx.AccountID(c.Request.Context())
	atype := string(contextx.AccountType(c.Request.Context()))
	if err := h.uc.Publish(c.Request.Context(), req.IDs, domainpkg.PublishParam{
		Status: "PUBLISHED", PublishAt: now,
		SenderAccountID: aid, SenderAccountType: atype, UpdatedBy: aid,
	}); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) revoke(c *gin.Context) {
	var req domainpkg.IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Revoke(c.Request.Context(), req.IDs, domainpkg.RevokeParam{
		Status: "REVOKED", RevokedAt: time.Now().UTC(),
	}); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) pin(c *gin.Context) {
	var req domainpkg.PinParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.Pin(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) portalList(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	// 门户公告列表固定查 ANNOUNCEMENT（对齐 voxel-boot portalList → pagePublished(..., KIND_ANNOUNCEMENT)）。
	if strings.TrimSpace(q.Kind) == "" {
		q.Kind = "ANNOUNCEMENT"
	}
	accountType, accountID := sessionScope(c)
	rows, total, current, size, err := h.uc.PagePublished(c.Request.Context(), q, accountType, accountID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (h *Handler) myPage(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	accountType, accountID := sessionScope(c)
	rows, total, current, size, err := h.uc.PagePublished(c.Request.Context(), q, accountType, accountID)
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
	accountType, accountID := sessionScope(c)
	row, err := h.uc.MyDetail(c.Request.Context(), q.ID, accountType, accountID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "notice not found")
		return
	}
	response.OK(c, row)
}

// sessionScope 从上下文取账户类型与 ID（匿名访问时类型取 PORTAL、ID 为空）。
func sessionScope(c *gin.Context) (string, string) {
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		return string(security.AccountPortal), ""
	}
	return string(sess.AccountType), sess.AccountID
}

func (h *Handler) unreadCount(c *gin.Context) {
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.OK(c, 0)
		return
	}
	total, err := h.uc.UnreadCount(c.Request.Context(), string(sess.AccountType), sess.AccountID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, total)
}

func (h *Handler) markRead(c *gin.Context) {
	var req domainpkg.ReadParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	now := time.Now().UTC()
	_ = h.uc.MarkReads(c.Request.Context(), string(sess.AccountType), sess.AccountID, req.IDs, now)
	response.OK(c, nil)
}

func (h *Handler) markAllRead(c *gin.Context) {
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	if err := h.uc.MarkAllRead(c.Request.Context(), string(sess.AccountType), sess.AccountID, time.Now().UTC()); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

// Package filehttp HTTP 触发器（xfg-ddd trigger）。
//
// Author: Charlie

//
// Author: Charlie

package filehttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/file"

	filecase "voxel-gin-admin/internal/cases/sys/file"

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
	uc *filecase.FileCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *filecase.FileCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/file/upload", admin, middleware.RequirePermission(d.Perms, "sys:file:upload", "文件上传"), middleware.OperationAudit(d.Audit, "sys_file", "upload"), h.upload)
		api.POST("/v1/admin/sys/file/delete", admin, middleware.RequirePermission(d.Perms, "sys:file:delete", "文件删除"), middleware.OperationAudit(d.Audit, "sys_file", "delete"), h.delete)
		api.POST("/v1/admin/sys/file/update", admin, middleware.RequirePermission(d.Perms, "sys:file:update", "文件更新"), middleware.OperationAudit(d.Audit, "sys_file", "update"), h.update)
		api.GET("/v1/admin/sys/file/detail", admin, middleware.RequirePermission(d.Perms, "sys:file:detail", "文件详情"), h.detail)
		api.GET("/v1/admin/sys/file/page", admin, middleware.RequirePermission(d.Perms, "sys:file:page", "文件分页"), h.page)
		api.POST("/v1/admin/sys/file/list_by_ids", admin, middleware.RequirePermission(d.Perms, "sys:file:detail", "文件批量查询"), h.listByIDs)
		api.GET("/v1/admin/sys/file/download", admin, middleware.RequirePermission(d.Perms, "sys:file:url", "文件下载"), h.download)
		api.POST("/v1/admin/sys/file/url", admin, middleware.RequirePermission(d.Perms, "sys:file:url", "文件URL"), h.url)
		api.POST("/v1/admin/sys/file/presigned_url", admin, middleware.RequirePermission(d.Perms, "sys:file:presignedurl", "文件预签名URL"), h.presignedURL)
	}
}

func (h *Handler) upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "file required")
		return
	}
	row, err := h.uc.Upload(c.Request.Context(), fh, c.PostForm("storage_provider"), accountID(c))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, row)
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
	response.OK(c, row)
}

func (h *Handler) page(c *gin.Context) {
	var q domainpkg.PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := h.uc.Page(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}

func (h *Handler) listByIDs(c *gin.Context) {
	var body domainpkg.IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	rows, err := h.uc.ListByIDs(c.Request.Context(), body.IDs)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) download(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, rc, err := h.uc.OpenByID(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	defer rc.Close()
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Disposition", domainpkg.ContentDispositionAttachment(row.OriginalName))
	c.DataFromReader(http.StatusOK, row.Size, "application/octet-stream", rc, nil)
}

func (h *Handler) url(c *gin.Context) {
	var req domainpkg.ObjectNameParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := h.uc.URL(c.Request.Context(), req.ObjectName)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, out)
}

func (h *Handler) presignedURL(c *gin.Context) {
	var req domainpkg.ObjectNameParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := h.uc.PresignedURL(c.Request.Context(), req.ObjectName)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, out)
}

// ---- Portal（仅本人文件，对齐 voxel-boot PortalFileController.assertOwnedByCurrent）----

func (h *Handler) portalUpload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "file required")
		return
	}
	row, err := h.uc.Upload(c.Request.Context(), fh, c.PostForm("storage_provider"), accountID(c))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) portalDetail(c *gin.Context) {
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
	if err := h.uc.AssertOwnedByCurrent(row, accountID(c)); err != nil {
		response.Fail(c, http.StatusForbidden, 403, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) portalListByIDs(c *gin.Context) {
	var body domainpkg.IDsParam
	if err := bind.JSON(c, &body); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	rows, err := h.uc.ListByIDs(c.Request.Context(), body.IDs)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	me := accountID(c)
	filtered := make([]domainpkg.File, 0, len(rows))
	for i := range rows {
		if h.uc.AssertOwnedByCurrent(&rows[i], me) == nil {
			filtered = append(filtered, rows[i])
		}
	}
	response.OK(c, filtered)
}

func (h *Handler) portalDownload(c *gin.Context) {
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
	if err := h.uc.AssertOwnedByCurrent(row, accountID(c)); err != nil {
		response.Fail(c, http.StatusForbidden, 403, err.Error())
		return
	}
	rc, err := h.uc.ProviderFor(c.Request.Context(), row).Get(c.Request.Context(), domainpkg.ToObjectKey(row.ObjectName))
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	defer rc.Close()
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Disposition", domainpkg.ContentDispositionAttachment(row.OriginalName))
	c.DataFromReader(http.StatusOK, row.Size, "application/octet-stream", rc, nil)
}

func (h *Handler) portalURL(c *gin.Context) {
	h.portalObjectName(c, func(ctx *gin.Context, objectName string) (*domainpkg.URLResult, error) {
		return h.uc.URL(ctx, objectName)
	})
}

func (h *Handler) portalPresignedURL(c *gin.Context) {
	h.portalObjectName(c, func(ctx *gin.Context, objectName string) (*domainpkg.URLResult, error) {
		return h.uc.PresignedURL(ctx, objectName)
	})
}

// portalObjectName 门户端按 object_name 获取 URL：先校验本人归属。
func (h *Handler) portalObjectName(c *gin.Context, fn func(ctx *gin.Context, objectName string) (*domainpkg.URLResult, error)) {
	var req domainpkg.ObjectNameParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	key := domainpkg.ToObjectKey(req.ObjectName)
	if key == "" {
		response.Fail(c, http.StatusNotFound, 404, "file not found")
		return
	}
	row, err := h.uc.Repo.FindByObjectName(c.Request.Context(), key)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "file not found")
		return
	}
	if err := h.uc.AssertOwnedByCurrent(row, accountID(c)); err != nil {
		response.Fail(c, http.StatusForbidden, 403, err.Error())
		return
	}
	out, err := fn(c, key)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.OK(c, out)
}

func accountID(c *gin.Context) string {
	if sess := contextx.Session(c.Request.Context()); sess != nil {
		return sess.AccountID
	}
	return ""
}

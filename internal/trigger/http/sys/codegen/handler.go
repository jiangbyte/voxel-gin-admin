// Package codegenhttp HTTP 触发器。
//
// Author: Charlie

//
// Author: Charlie

package codegenhttp

import (
	domainpkg "voxel-gin-admin/internal/domain/sys/codegen"

	codegencase "voxel-gin-admin/internal/cases/sys/codegen"

	"net/http"

	"github.com/gin-gonic/gin"

	"voxel-gin-admin/internal/infrastructure/core/bind"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/types/schema"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
	"voxel-gin-admin/internal/infrastructure/platform/module"
)

// Handler HTTP 适配：只调用 Case。
type Handler struct {
	uc *codegencase.CodegenCase
}

// NewHandler 构造 Handler。
func NewHandler(uc *codegencase.CodegenCase) *Handler {
	return &Handler{uc: uc}
}


func (h *Handler) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		api.POST("/v1/admin/sys/codegen/create", admin, middleware.RequirePermission(d.Perms, "sys:codegen:create", "代码生成创建"), middleware.OperationAudit(d.Audit, "sys_codegen", "create"), h.create)
		api.POST("/v1/admin/sys/codegen/update", admin, middleware.RequirePermission(d.Perms, "sys:codegen:update", "代码生成更新"), middleware.OperationAudit(d.Audit, "sys_codegen", "update"), h.update)
		api.POST("/v1/admin/sys/codegen/delete", admin, middleware.RequirePermission(d.Perms, "sys:codegen:delete", "代码生成删除"), middleware.OperationAudit(d.Audit, "sys_codegen", "delete"), h.delete)
		api.GET("/v1/admin/sys/codegen/detail", admin, middleware.RequirePermission(d.Perms, "sys:codegen:detail", "代码生成详情"), h.detail)
		api.GET("/v1/admin/sys/codegen/page", admin, middleware.RequirePermission(d.Perms, "sys:codegen:page", "代码生成分页"), h.page)
		api.GET("/v1/admin/sys/codegen/tables", admin, middleware.RequirePermission(d.Perms, "sys:codegen:tables", "代码生成表列表"), h.tables)
		api.GET("/v1/admin/sys/codegen/table-columns", admin, middleware.RequirePermission(d.Perms, "sys:codegen:tables", "代码生成表列元数据"), h.tableColumns)
		api.GET("/v1/admin/sys/codegen/fields", admin, middleware.RequirePermission(d.Perms, "sys:codegen:detail", "代码生成字段"), h.fields)
		api.POST("/v1/admin/sys/codegen/fields/update-batch", admin, middleware.RequirePermission(d.Perms, "sys:codegen:update", "代码生成字段批量更新"), middleware.OperationAudit(d.Audit, "sys_codegen", "update"), h.updateFieldsBatch)
		api.GET("/v1/admin/sys/codegen/parent-resources", admin, middleware.RequirePermission(d.Perms, "sys:codegen:detail", "代码生成父级资源"), h.parentResources)
		api.GET("/v1/admin/sys/codegen/preview", admin, middleware.RequirePermission(d.Perms, "sys:codegen:preview", "代码生成预览"), h.preview)
		api.GET("/v1/admin/sys/codegen/download", admin, middleware.RequirePermission(d.Perms, "sys:codegen:download", "代码生成下载"), h.download)
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

func (h *Handler) tables(c *gin.Context) {
	rows, err := h.uc.Tables(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) tableColumns(c *gin.Context) {
	tableName := c.Query("table_name")
	if tableName == "" {
		response.Fail(c, http.StatusBadRequest, 400, "table_name required")
		return
	}
	rows, err := h.uc.TableColumns(c.Request.Context(), tableName)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) fields(c *gin.Context) {
	var q domainpkg.FieldQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	rows, err := h.uc.Fields(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) updateFieldsBatch(c *gin.Context) {
	var req domainpkg.FieldsUpdateBatchParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := h.uc.UpdateFieldsBatch(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) parentResources(c *gin.Context) {
	moduleID := c.Query("module_id")
	rows, err := h.uc.ParentResources(c.Request.Context(), moduleID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handler) preview(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := h.uc.Preview(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, out)
}

func (h *Handler) download(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	raw, name, err := h.uc.DownloadZip(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	c.Header("Content-Disposition", "attachment; filename=\""+name+"\"")
	c.Data(http.StatusOK, "application/zip", raw)
}

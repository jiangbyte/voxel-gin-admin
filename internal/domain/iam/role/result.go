// internal/modules/iam/role/result.go 出参定义。
//
// Author: Charlie

package role

import (
	"voxel-gin-admin/internal/domain/iam/client"
	"voxel-gin-admin/internal/domain/iam/relation"
	"voxel-gin-admin/internal/domain/iam/resource"
)

// OwnResourceResult 角色已拥有管理端资源授权结果。
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult 角色已拥有客户端资源授权结果。
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

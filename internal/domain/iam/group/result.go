// internal/modules/iam/group/result.go 出参定义。
//
// Author: Charlie

package group

import (
	"voxel-gin-admin/internal/domain/iam/client"
	"voxel-gin-admin/internal/domain/iam/relation"
	"voxel-gin-admin/internal/domain/iam/resource"
	"voxel-gin-admin/internal/domain/iam/role"
)

// OwnRoleResult 用户组已拥有角色结果。
//
// Author: Charlie
type OwnRoleResult struct {
	ID      string      `json:"id"`
	Roles   []role.Role `json:"roles"`
	RoleIDs []string    `json:"role_ids"`
}

// OwnResourceResult 用户组已拥有管理端资源授权结果。
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult 用户组已拥有客户端资源授权结果。
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

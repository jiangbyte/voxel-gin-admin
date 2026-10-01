// Package all 汇总 blank import 全部内置业务模块，触发 init 自注册。
//
// 二次开发推荐：整仓 Git 合并上游；cmd 保留对本包的 blank import，
// 自有扩展在自己的 main 里额外 _ import，减少与官方 all 的合并冲突。
// 复杂场景直接改本仓 framework/（非外部依赖升级模型）。
//
// Author: Charlie
package all

import (
	_ "voxel-gin-admin/internal/trigger/http/auth"
	_ "voxel-gin-admin/internal/trigger/http/sys/workspace"
	_ "voxel-gin-admin/internal/domain/internal/health"
	_ "voxel-gin-admin/internal/trigger/http/iam/account"
	_ "voxel-gin-admin/internal/trigger/http/iam/client"
	_ "voxel-gin-admin/internal/trigger/http/iam/dept"
	_ "voxel-gin-admin/internal/trigger/http/iam/group"
	_ "voxel-gin-admin/internal/trigger/http/iam/position"
	_ "voxel-gin-admin/internal/trigger/http/iam/relation"
	_ "voxel-gin-admin/internal/trigger/http/iam/resource"
	_ "voxel-gin-admin/internal/trigger/http/iam/role"
	_ "voxel-gin-admin/internal/domain/profile/admin"
	_ "voxel-gin-admin/internal/trigger/http/profile/identity"
	_ "voxel-gin-admin/internal/domain/profile/portal"
	_ "voxel-gin-admin/internal/trigger/http/sys/audit"
	_ "voxel-gin-admin/internal/trigger/http/sys/banner"
	_ "voxel-gin-admin/internal/trigger/http/sys/codegen"
	_ "voxel-gin-admin/internal/trigger/http/sys/config"
	_ "voxel-gin-admin/internal/trigger/http/sys/dict"
	_ "voxel-gin-admin/internal/trigger/http/sys/feedback"
	_ "voxel-gin-admin/internal/trigger/http/sys/file"
	_ "voxel-gin-admin/internal/trigger/http/sys/job"
	_ "voxel-gin-admin/internal/trigger/http/sys/notice"
	_ "voxel-gin-admin/internal/trigger/http/sys/weak_password"
)

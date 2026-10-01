# Voxel Gin Admin

管理端 Go 后端（Gin），xfg-ddd 布局：

- `internal/types` — 共享类型（errors / schema / response）
- `internal/api` — API 契约
- `internal/domain` — 业务域
- `internal/infrastructure` — 框架与中间件
- `internal/cases` — 用例编排
- `internal/trigger` — HTTP 等入口
- `cmd/app` — 进程入口

仅挂载 `/api/v1/admin/**`（+ internal）。默认端口 `8200`。

```bash
go build -o bin/voxel-gin-admin ./cmd/app
./bin/voxel-gin-admin
```

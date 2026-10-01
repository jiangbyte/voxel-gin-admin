# Voxel Gin Admin

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Supported-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-4169E1?logo=postgresql&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-Supported-4479A1?logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Supported-DC382D?logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)
![Version](https://img.shields.io/badge/version-1.1.0--SNAPSHOT-orange)

**Voxel Gin Admin** 是面向中后台管理端的 Go / Gin 后端：DDD 布局，仅挂载 **Admin** API（`/api/v1/admin/**`，含 internal）。

## 目录

- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [工程结构](#工程结构)
- [快速开始](#快速开始)
- [默认账号](#默认账号)
- [License](#license)

## 功能特性

API 前缀为 `/api/v1/admin/*`：

| 模块 | 说明 |
| --- | --- |
| 账号体系 | ADMIN 会话（Redis 不透明 Token）；密码 RSA、验证码、失败锁定与限流；可配置三方 OAuth |
| RBAC 权限 | 账号 / 角色 / 部门 / 用户组 / 岗位；菜单、按钮与 API 资源授权；在线会话踢出 |
| 系统管理 | 字典、动态配置（`sys_config`，敏感项可加密）、Banner、公告 / 通知、意见反馈、弱口令库 |
| 对象存储 | S3 兼容存储，引擎与凭证走运行时配置 |
| 运维能力 | 操作审计、登录日志、工作台概览、内置任务调度（`sys_job`；DB 扫描 + Redis 锁 + cron） |
| 代码生成 | 单表 / 树表 / 主子表，预览与 ZIP 下载 |
| 实名认证 | 工单提交与审核、敏感字段加密存储 |
| 业务扩展 | 领域模块可按同样 DDD 模式横向扩展 |

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 后端 | Go 1.25+ · Gin · 单 module（`cmd/app`） |
| 持久化 | PostgreSQL / MySQL · GORM |
| 缓存 / 会话 | Redis（go-redis）· 不透明会话 Token |
| 配置 | Viper（`config.yaml`）+ 运行时 `sys_config` |
| 文档 | swaggo · OpenAPI `/v3/api-docs` · Swagger UI `/swagger-ui/index.html` |
| 其他 | AWS SDK v2（S3）· zap · robfig/cron · snowflake |

## 工程结构

```text
voxel-gin-admin/
├── cmd/app/                    # 可启动入口
├── internal/
│   ├── types/                  # 共享类型（errors / schema / response）
│   ├── api/                    # API 契约
│   ├── domain/                 # 业务域（auth / iam / sys / profile …）
│   ├── infrastructure/         # 框架与中间件
│   ├── cases/                  # 用例编排
│   ├── trigger/                # HTTP 等入口
│   └── app/                    # 装配根
├── configs/config.example.yaml
└── scripts/voxel_gin.sql
```

| 文件 | 用途 |
| --- | --- |
| `scripts/voxel_gin.sql` | MySQL 全量建表、种子与表/列 `COMMENT`（`sys_job.handler` 为 Gin 原生 key） |

## 快速开始

### 环境要求

- Go **1.25+**
- MySQL 8+、Redis

### 1. 初始化数据库

```bash
mysql -u root -p -e "CREATE DATABASE voxel_gin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p voxel_gin < scripts/voxel_gin.sql
```

复制配置并按需修改：

```bash
cp configs/config.example.yaml config.yaml
# 编辑 db / redis / app 等；默认端口 8200，库名 voxel_gin
```

### 2. 启动后端

```bash
go run ./cmd/app
# 或
go build -o bin/voxel-gin-admin ./cmd/app && ./bin/voxel-gin-admin
```

| 项 | 地址 |
| --- | --- |
| API | http://127.0.0.1:8200 |
| Swagger UI | http://127.0.0.1:8200/swagger-ui/index.html |
| OpenAPI JSON | http://127.0.0.1:8200/v3/api-docs |

## 默认账号

| 端 | 地址 | 账号 | 密码 | 说明 |
| --- | --- | --- | --- | --- |
| Admin | http://127.0.0.1:8200 | `superadmin` | `123456` | 超级管理员（`*:*:*`） |

仅供本地演示。部署后请修改默认密码与敏感配置。更多种子见 `scripts/voxel_gin.sql`。

## License

本项目基于 [Apache License 2.0](LICENSE) 开源。完整条款见 [LICENSE](LICENSE)，版权声明见 [NOTICE](NOTICE)。

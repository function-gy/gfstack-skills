---
name: gfstack-overview
description: "GoFrame v2 enterprise architecture overview. Covers project directory layout, HTTP request flow, and key constraint checklist. USE FOR: understanding project structure, request lifecycle, architecture constraints, anti-pattern checks. DO NOT USE FOR: writing specific code — load gfstack-api/gfstack-logic/gfstack-data instead."
---

# gfstack-overview

## 1. Project Directory Layout

TRIGGER: project structure, directory layout, module organization, new project, project scaffolding

```
server/
├── main.go                     # Entry: global.Init(ctx) → cmd.Main.Run(ctx)
├── go.mod                      # Module: hotgo (Go 1.24+, GF v2.9+)
├── api/                        # API contract layer (Req/Res structs with g.Meta)
│   └── {module}/               # e.g. admin/, api/
│       ├── {module}.go         # Controller interface: I{Module}V1
│       └── v1/                 # Per-entity route definitions: {entity}.go
├── internal/
│   ├── cmd/                    # CLI commands: http, cron, queue, tools
│   ├── consts/                 # Application constants & error codes
│   ├── controller/             # Thin handlers (call services only, NO business logic)
│   │   └── {module}/           # e.g. admin/admin/, admin/api/
│   ├── dao/                    # Generated DAO + custom extensions
│   │   └── internal/           # Generated DAO internals (table, columns, DB())
│   ├── global/                 # System bootstrap (config, cache, trace, cluster sync)
│   ├── library/                # Reusable infrastructure (cache, casbin, hgorm, etc.)
│   ├── logic/                  # Business logic implementations
│   │   ├── logic.go            # Import aggregator (blank imports all logic packages)
│   │   └── {domain}/           # e.g. sys/, admin/, middleware/, hook/
│   ├── model/                  # Data models
│   │   ├── entity/             # ORM entity structs (1:1 DB table mapping, auto-generated)
│   │   ├── do/                 # Data Objects for ORM write operations (auto-generated)
│   │   ├── input/              # Input/output DTOs
│   │   │   ├── form/           # Base types: PageReq, PageRes, Sorter, SwitchReq
│   │   │   └── {domain}in/     # Domain inputs: sysin/, adminin/, apiin/
│   │   ├── config.go           # Configuration model structs
│   │   ├── response.go         # Standard HTTP Response envelope
│   │   └── context.go          # Request context structs (Identity, Context)
│   ├── packed/                 # Resource packing (gres)
│   ├── router/                 # Route registration (admin.go, api.go)
│   └── service/                # Service interfaces + accessor/register functions
├── utility/                    # Generic utility functions
├── cron/                       # Cron job implementations
├── queue/                      # Queue consumer implementations
├── manifest/config/            # YAML config files
├── resource/                   # Static resources (templates, i18n)
├── storage/                    # Runtime storage (cache, certs, generated SQL)
└── logs/                       # Log output directories
```

---

## 2. Request Flow

TRIGGER: request lifecycle, middleware order, HTTP pipeline

```
HTTP Request
  → Global Middleware (RequestLog → Ctx → CORS → Blacklist → DemoLimit → ResponseHandler)
    → Route Group Middleware (AdminAuth / Develop)
      → Controller (thin handler — service calls only)
        → Service Interface Call
          → Logic Implementation (business logic)
            → DAO / gdb.Model (ORM)
            → Entity / DO / Input DTOs (data models)
```

---

## 3. Key Constraints (索引清单)

> 以下为前文规则的快速索引，详细说明和代码示例见对应子技能。

| # | 约束 | 详见 |
|---|------|------|
| 1 | Controllers MUST be thin — NO business logic | gfstack-api §3 |
| 2 | Logic packages MUST NOT import other logic packages — use service interfaces | gfstack-logic |
| 3 | DAO / Entity / DO are auto-generated — NEVER manually create | gfstack-data §1~§3 |
| 4 | Prefer DO structs over raw `g.Map` for DB operations; when `g.Map` is necessary, keys MUST use `dao.Xxx.Columns().FieldName` | gfstack-data §5 |
| 5 | DO NOT manually set `created_at` / `updated_at` / `deleted_at` | gfstack-data §1 |
| 6 | DO NOT manually add `WhereNull(deleted_at)` — framework auto-filters soft delete | gfstack-data §1 |
| 7 | Logic layer errors MUST use `gerror.Wrap()` | gfstack-style §1 |
| 8 | List responses MUST return empty slice, never nil | gfstack-api §1 Rules |
| 9 | ORM operations MUST use `Fields()` with a struct for field whitelisting | gfstack-data §4 |
| 10 | Entity structs use `orm` + `json` tags; Input DTOs use only `json` + `dc` tags | gfstack-style §3 |
| 11 | Error codes MUST be defined in `internal/consts/` using `gcode.New()` | gfstack-style §1 |
| 12 | `:=` is FORBIDDEN except in `for i := 0` loop initializers — use `var` blocks | gfstack-style §3 |
| 13 | Controller and logic layers MUST NOT contain unused functions, variables, or constants | gfstack-api, gfstack-logic |

### 子技能速查

| 技能 | 覆盖章节 | 何时加载 |
|------|---------|----------|
| gfstack-api | API 定义、Controller 接口、Controller 实现 | 写 API 接口或 Controller |
| gfstack-logic | Service 接口、Logic 实现、逻辑层约束 | 写业务逻辑 |
| gfstack-data | Entity/DO/DTO/DAO/ORM | 数据建模或数据库操作 |
| gfstack-route | 路由注册、中间件 | 配置路由或中间件 |
| gfstack-infra | Token、服务启动、响应格式 | 基础设施相关 |
| gfstack-style | 错误码、校验、命名、编码风格 | 代码规范检查 |

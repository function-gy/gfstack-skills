---
name: gfstack
description: "GoFrame v2 enterprise architecture skill index. Triggers for general architecture questions — loads gfstack-* sub-skills for specific implementations. USE FOR: project structure overview, understanding architecture. For specific code: use gfstack-api, gfstack-logic, gfstack-data, gfstack-route, gfstack-infra, gfstack-style."
---

# HotGo Architecture (模块索引)

> 本技能已拆分为 7 个独立子技能，按需加载，减少 token 浪费。

## 子技能速查

| 技能 | 覆盖内容 | 何时加载 |
|------|---------|----------|
| `gfstack-overview` | 项目目录布局、请求全链路、所有约束清单 | 理解架构全局 |
| `gfstack-api` | API 定义（g.Meta Req/Res）、Controller 接口、Controller 实现 | 写或修改 API 层 |
| `gfstack-logic` | Service 接口 + Logic 实现（header pattern、方法签名、init 注册） | 写业务逻辑 |
| `gfstack-data` | Entity / DO / Input DTO / Form / DAO / ORM 全模式 | 操作数据层 |
| `gfstack-route` | 路由注册（group.Bind）、中间件 | 配置路由/中间件 |
| `gfstack-infra` | Token 系统、服务启动（main.go）、响应格式 | 基础设施 |
| `gfstack-style` | 错误码、校验、变量声明、命名返回值、godoc、命名约定 | 编码规范检查 |

## 规则冲突

如果多个子技能之间的规则冲突，以 **gfstack-style** 和 **gfstack-overview** 为准。

## 加载策略

- 写 API 接口：加载 `gfstack-api` + `gfstack-logic` + `gfstack-data`
- 加中间件：仅加载 `gfstack-route`
- 定义数据模型：仅加载 `gfstack-data`
- 代码审查：加载 `gfstack-style`

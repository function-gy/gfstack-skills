# gfstack — 企业级开发技能集

GoFrame v2 企业级项目架构知识库，已拆分为 7 个独立子技能，按需加载减少 token 浪费。

## 模块结构

```
skills/
├── hotgo-architecture/          # 入口索引（本目录）
│   ├── SKILL.md                 # 索引，指向各子技能
│   ├── README.md                # 本文件
│   └── examples/                # 代码示例
├── gfstack-overview/            # 架构总览：目录 + 请求流 + 约束清单
├── gfstack-api/                 # API 层：g.Meta Req/Res + Controller
├── gfstack-logic/               # 逻辑层：Service 接口 + Logic 实现
├── gfstack-data/                # 数据层：Entity/DO/DTO + DAO + ORM
├── gfstack-route/               # 路由层：Router + Middleware
├── gfstack-infra/               # 基础设施：Token + Bootstrap + Response
└── gfstack-style/               # 规范层：错误码 + 校验 + 命名 + 编码风格
```

## 按需加载策略

| 任务 | 需要加载的技能 |
|------|---------------|
| 写一个 CRUD API | gfstack-api + gfstack-logic + gfstack-data |
| 加/改中间件 | gfstack-route |
| 建数据表模型 | gfstack-data |
| 项目初始化/全局架构理解 | gfstack-overview |
| 代码审查/风格统一 | gfstack-style |
| Token 认证相关 | gfstack-infra |

## 规则优先级

子技能之间规则冲突时：**gfstack-style > gfstack-overview > 其他子技能**

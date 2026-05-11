# gfstack — GoFrame v2 企业级开发技能集

AI 辅助开发技能集，适用于 opencode、Claude Code、Codex 等 AI 编程助手，提供强制性的分层架构规范和编码风格约束。

## 快速安装

```bash
npx degit function-gy/gfstack-skills ~/.opencode/skills
```

或使用 git：

```bash
git clone git@github.com:function-gy/gfstack-skills.git ~/.opencode/skills
```

更新：

```bash
cd ~/.opencode/skills && git pull
```

安装后在 opencode 中输入 `/skills` 即可查看所有技能。

## 目录结构

```
skills/
├── README.md
├── gfstack/              # 入口索引
│   └── examples/         # 20 个代码示例
├── gfstack-overview/     # 架构总览：目录布局、请求流、约束清单
├── gfstack-api/          # API 层：g.Meta Req/Res + Controller
├── gfstack-logic/        # 逻辑层：Service 接口 + Logic 实现
├── gfstack-data/         # 数据层：Entity/DO/DTO + DAO + ORM
├── gfstack-route/        # 路由层：Router + Middleware
├── gfstack-infra/        # 基础设施：Token + 启动 + 响应格式
├── gfstack-style/        # 规范层：错误码 + 校验 + 命名 + 编码风格
├── gfstack-audit/        # 审计层：安全审查 + 代码检查清单
├── ui-ux-pro-max/        # UI/UX 设计规范
├── vue-best-practices/   # Vue.js 最佳实践
└── upgrade-skills/       # 技能自动升级工具
```

## 按需加载

gfstack 技能设计为按需加载，减少 token 消耗：

| 任务 | 加载的技能 |
|------|---------------|
| 写 CRUD API | gfstack-api + gfstack-logic + gfstack-data |
| 加中间件 | gfstack-route |
| 定义数据模型 | gfstack-data |
| 理解整体架构 | gfstack-overview |
| 代码审查 | gfstack-style + gfstack-audit |
| Token/认证 | gfstack-infra |

## 核心原则

1. 严格分层：Controller → Service 接口 → Logic 实现 → DAO
2. 接口驱动：`ISysXxx` 接口 + `RegisterSysXxx()` 注册模式
3. 自动生成：DAO / Entity / DO 由 `gf gen dao` 生成，禁止手动修改
4. 命名约定：`I` 前缀（接口）、`s` 前缀（Logic 结构体）、`New`（构造函数）、`Register`（注册函数）
5. 错误处理：`gerror.Wrap()` 包装 + `gcode.New()` 错误码
6. 编码风格：`:=` 仅用于 `for` 循环，变量用 `var` 块声明，函数用命名返回值，godoc 用中文注释

## License

MIT

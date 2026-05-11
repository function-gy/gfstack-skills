# gfstack — GoFrame v2 Enterprise Dev Skills

AI-assisted development skills based on [HotGo](https://github.com/bufanyun/hotgo) architecture for [opencode](https://github.com/anomalyco/opencode), enforcing strict layered architecture and coding conventions.

## Quick Install

```bash
git clone git@github.com:function-gy/gfstack-skills.git ~/.opencode/skills
```

To update:

```bash
cd ~/.opencode/skills && git pull
```

Run `/skills` in opencode to see all installed skills.

## Structure

```
skills/
├── README.md
├── gfstack/              # Entry index (this entry)
│   └── examples/         # 20 code examples
├── gfstack-overview/     # Architecture overview: layout, request flow, constraints
├── gfstack-api/          # API layer: g.Meta Req/Res + Controller
├── gfstack-logic/        # Logic layer: Service interface + Logic implementation
├── gfstack-data/         # Data layer: Entity/DO/DTO + DAO + ORM
├── gfstack-route/        # Route layer: Router + Middleware
├── gfstack-infra/        # Infrastructure: Token + Bootstrap + Response
├── gfstack-style/        # Standards: Error codes + Validation + Naming + Style
├── ui-ux-pro-max/        # UI/UX design guidelines
├── vue-best-practices/   # Vue.js best practices
└── upgrade-skills/       # Skill auto-upgrade utility
```

## On-Demand Loading

gfstack skills are designed for minimal token usage — only relevant skills load per task:

| Task | Skills Loaded |
|------|---------------|
| Write CRUD API | gfstack-api + gfstack-logic + gfstack-data |
| Add middleware | gfstack-route |
| Define data models | gfstack-data |
| Understand architecture | gfstack-overview |
| Code review | gfstack-style |
| Token/auth | gfstack-infra |

## Core Principles

1. Strict layering: Controller → Service → Logic → DAO
2. Interface-driven: `ISysXxx` interface + `RegisterSysXxx()` pattern
3. Auto-generated: DAO / Entity / DO by `gf gen dao`, never manually edit
4. Naming: `I` prefix (interface), `s` prefix (logic struct), `New` (constructor), `Register` (registration)
5. Errors: `gerror.Wrap()` + `gcode.New()` error codes
6. Style: `:=` only in `for` loops, `var` blocks, named return values, Chinese godoc comments

## License

MIT

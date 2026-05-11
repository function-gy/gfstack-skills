---
name: gfstack
description: "GoFrame v2 enterprise architecture skill index. Triggers for general architecture questions — loads gfstack-* sub-skills for specific implementations. USE FOR: project structure overview, understanding architecture, code review. For specific code: use gfstack-api, gfstack-logic, gfstack-data, gfstack-route, gfstack-infra, gfstack-style, gfstack-audit."
---

# gfstack — Skill Index

> Split into 7 independent sub-skills for on-demand loading, minimizing token waste.

## Sub-Skills

| Skill | Coverage | When to Load |
|-------|----------|--------------|
| `gfstack-overview` | Directory layout, request flow, all constraint checklist | Understanding overall architecture |
| `gfstack-api` | API definitions (g.Meta Req/Res), controller interface, controller implementation | Writing or modifying API layer |
| `gfstack-logic` | Service interface + logic implementation (header pattern, method signature, init registration) | Writing business logic |
| `gfstack-data` | Entity / DO / Input DTO / Form / DAO / ORM patterns | Data layer operations |
| `gfstack-route` | Route registration (group.Bind), middleware | Configuring routes/middleware |
| `gfstack-infra` | Token system, bootstrap (main.go), response format | Infrastructure tasks |
| `gfstack-style` | Error codes, validation, variable declaration, named returns, godoc, naming conventions | Code style reviews |
| `gfstack-audit` | Security audit checklist: injection, sensitive data, concurrency, panics, stability | Code security review |

## Rule Precedence

In case of conflicts between sub-skills: **gfstack-style** > **gfstack-overview** > others.

## Loading Strategy

- Writing a CRUD API: `gfstack-api` + `gfstack-logic` + `gfstack-data`
- Adding middleware: `gfstack-route` only
- Defining data models: `gfstack-data` only
- Code review: `gfstack-style` + `gfstack-audit`
- Security audit: `gfstack-audit` only

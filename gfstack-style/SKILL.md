---
name: gfstack-style
description: "HotGo coding standards. Covers error code specification, validation rules, variable declaration conventions, named return values, godoc comments, struct field tags (json+dc), inline if error handling, and naming conventions. TRIGGER: error code, gcode.New, error handling, errcode, validation, Filter, g.Validator, PreFilter, v tag, variable declaration, var block, short variable declaration, named return, coding style, naming convention, I prefix, s prefix, New function, Register function. DO NOT TRIGGER: API definitions (see gfstack-api), data models (see gfstack-data), middleware (see gfstack-route)."
---

# gfstack-style

## 1. Error Code Specification

TRIGGER: error code, gcode.New, error handling, errcode, CodeNotAuthorized

All business error codes defined in `internal/consts/` using `gcode.New(code, message, detail)`:

| Range | Module |
|-------|--------|
| 0 | Success |
| 1xxxx | Common errors (param, auth, permission) |
| 2xxxx | User module |
| 3xxxx | Dealer module |
| ... | Increment by module |

Usage: `gerror.NewCode(consts.CodeDealerNotFound, "经销商不存在")` → `gerror.Wrap(err, "..." )`

## 2. Validation

TRIGGER: validation, Filter, g.Validator, PreFilter, struct tag validation, v tag

```go
Id int64 `json:"id" v:"required#id不能为空" dc:"id"`

func (in *FooEditInp) Filter(ctx context.Context) (err error) {
    if err = g.Validator().Rules("required").Data(in.Name).Messages("名称不能为空").Run(ctx); err != nil {
        return err.Current()
    }
    return
}
```

> Filter 方法完整定义和 Input DTO 结构见 `gfstack-data` §3。

---

## 3. Variable Declaration Convention

TRIGGER: variable declaration, var block, short variable declaration, named return, coding style

- **ALL comments MUST be in Chinese** — code comments, function doc comments, inline explanations
- Function-scope variables: `var` block at top or named return values
- `:=` is FORBIDDEN except in `for i := 0` loop initializers
- Multiple variables/constants use batch declaration blocks

<!-- PATCHED: 2026-05-09 -->
**Multiple variables MUST be grouped into a single `var ()` block.** Do NOT write individual `var` statements:

```go
// WRONG — individual var statements
var totalUsers int
var submittedUsers int
var lateSubmitUsers int
var totalHours float64

// CORRECT — single var () block
var (
    totalUsers      int
    submittedUsers  int
    lateSubmitUsers int
    totalHours      float64
)
```

<!-- PATCHED: 2026-05-09 -->
**Every function/method MUST have a godoc comment** above the `func` declaration describing its purpose. Be detailed but concise — do not exceed 6 lines. Format:

```go
// FooList 获取Foo分页列表。
// 根据传入的筛选条件（名称、状态）查询Foo表，
// 返回分页后的数据列表和总记录数。
// 列表默认按 Id DESC 排序。
func (s *sSysFoo) FooList(ctx context.Context, in *sysin.FooListInp) (out *sysin.FooListModel, err error) {
    ...
}
```

Key rules:
- Starts with `// {FuncName}` followed by Chinese description
- Explains what the function does, key input/output, and notable side effects (DB ops, cache, etc.)
- 2–6 lines total; never a bare one-liner with no detail
- Applies to all functions: logic methods, controller handlers, utility functions, middleware

<!-- PATCHED: 2026-05-09 -->
**Every struct field MUST have `json` and `dc` tags.** No field is allowed without both metadata tags:

```go
// WRONG — missing dc tag
type PageRes struct {
    Page      int   `json:"page"`
    PageSize  int   `json:"page_size"`
}

// CORRECT — json + dc on every field
type PageRes struct {
    Page      int   `json:"page"      dc:"页码"`
    PageSize  int   `json:"page_size" dc:"每页条数"`
    PageCount int   `json:"page_count" dc:"总页数"`
    Total     int64 `json:"total"     dc:"总数"`
}
```

This applies to ALL structs: API Req/Res, Entity, DO, Input DTOs, form types, config structs, and any custom struct. Entity structs additionally require `orm` tag (see `gfstack-data` §1).

<!-- PATCHED: 2026-05-10 -->
**Function return types MUST use named return values.** The function body should use bare `return` instead of explicit `return x, y`:

```go
// WRONG — unnamed return values, explicit return
func (s *sSysFoo) List(ctx context.Context, in *sysin.FooListInp) ([]*sysin.FooListModel, int, error) {
    ...
    return list, totalCount, nil
}

// CORRECT — named return values, bare return
func (s *sSysFoo) List(ctx context.Context, in *sysin.FooListInp) (list []*sysin.FooListModel, totalCount int, err error) {
    ...
    return
}
```

This applies to ALL functions: logic methods, controller handlers, service interface definitions, and utility functions. Only exception: empty return types (e.g. `func foo() {}`).

<!-- PATCHED: 2026-05-10 -->
**Error-returning function calls MUST use inline `if` assignment.** Do NOT separate variable assignment from error check:

```go
// WRONG — assignment and error check on separate lines
wo, err = s.getByID(ctx, workOrderID)
if err != nil {
    return
}

// CORRECT — inline if with semicolon
if wo, err = s.getByID(ctx, workOrderID); err != nil {
    return
}
```

This pattern ensures the error is checked immediately at the call site and prevents accidental use of the variable before the error is handled.

---

## 4. Naming Conventions

TRIGGER: naming convention, I prefix, s prefix, New function, Register function

| Pattern | Convention |
|---|---|
| Service interface | `I{Module}{Domain}` e.g. `ISysDealer` |
| Logic struct | `s{Module}{Domain}` e.g. `sSysDealer` |
| Accessor | `{Module}{Domain}()` e.g. `SysDealer()` |
| Register | `Register{Module}{Domain}()` e.g. `RegisterSysDealer()` |
| Constructor | `New{Module}{Domain}()` e.g. `NewSysDealer()` |
| Controller struct | `ControllerV1` (empty) |
| API Req/Res | `{Entity}{Action}Req/Res` |
| Input DTO | `{Entity}{Action}Inp` |
| Input Model | `{Entity}{Action}Model` |
| Field filter | `{Entity}UpdateFields` / `{Entity}InsertFields` |
| DAO variable | `dao.{Entity}` |

> Service 接口和 Logic 实现模式见 `gfstack-logic`。Controller 实现见 `gfstack-api` §3。DAO 变量命名见 `gfstack-data` §5。

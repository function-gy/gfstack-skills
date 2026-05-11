---
name: gfstack-logic
description: "HotGo service and logic layer specification. Covers service interface pattern, logic implementation with header pattern, method signatures, and init registration. TRIGGER: service interface, ISys, Register, accessor function, logic implementation, sSys, init register, business logic, service locator pattern. DO NOT TRIGGER: API definitions (see gfstack-api), ORM operations (see gfstack-data)."
---

# gfstack-logic

## 1. Service Interface + Implementation Pattern

### 1.1 Service Interface

TRIGGER: service interface, ISys, service locator, Register, accessor function, ISysFoo

Location: `internal/service/{domain}.go`

```go
package service

type (
    ISysFoo interface {
        Model(ctx context.Context, option ...*handler.Option) *gdb.Model
        List(ctx context.Context, in *sysin.FooListInp) (list []*sysin.FooListModel, totalCount int, err error)
        Edit(ctx context.Context, in *sysin.FooEditInp) (err error)
        View(ctx context.Context, in *sysin.FooViewInp) (res *sysin.FooViewModel, err error)
        Delete(ctx context.Context, in *sysin.FooDeleteInp) (err error)
        Switch(ctx context.Context, in *sysin.FooSwitchInp) (err error)
    }
)

var localSysFoo ISysFoo

func SysFoo() ISysFoo {
    if localSysFoo == nil {
        panic("implement not found for interface ISysFoo, forgot register?")
    }
    return localSysFoo
}

func RegisterSysFoo(i ISysFoo) {
    localSysFoo = i
}
```

### 1.2 Logic Implementation

TRIGGER: logic implementation, sSys, init register, business logic, NewSysFoo

Location: `internal/logic/{domain}/{entity}.go`

```go
package sys

type sSysFoo struct{}

func NewSysFoo() *sSysFoo {
    return &sSysFoo{}
}

func init() {
    service.RegisterSysFoo(NewSysFoo())
}

func (s *sSysFoo) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
    return handler.Model(dao.Foo.Ctx(ctx), option...)
}

func (s *sSysFoo) List(ctx context.Context, in *sysin.FooListInp) (list []*sysin.FooListModel, totalCount int, err error) {
    mod := s.Model(ctx)
    mod = mod.Fields(sysin.FooListModel{})

    if in.Name != "" {
        mod = mod.WhereLike(dao.Foo.Columns().Name, "%"+in.Name+"%")
    }

    mod = mod.Page(in.Page, in.PerPage)
    mod = mod.OrderDesc(dao.Foo.Columns().Id)

    if err = mod.ScanAndCount(&list, &totalCount, false); err != nil {
        err = gerror.Wrap(err, "获取列表失败，请稍后重试！")
        return
    }
    return
}

func (s *sSysFoo) Edit(ctx context.Context, in *sysin.FooEditInp) (err error) {
    if err = validate.PreFilter(ctx, in); err != nil {
        return
    }
    if err = hgorm.IsUnique(ctx, &dao.Foo, g.Map{dao.Foo.Columns().Name: in.Name}, "名称已存在", in.Id); err != nil {
        return
    }
    return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
        if in.Id > 0 {
            if _, err = s.Model(ctx).Fields(sysin.FooUpdateFields{}).WherePri(in.Id).Data(in).Update(); err != nil {
                err = gerror.Wrap(err, "修改失败，请稍后重试！")
            }
            return
        }
        if _, err = s.Model(ctx, &handler.Option{FilterAuth: false}).
            Fields(sysin.FooInsertFields{}).Data(in).OmitEmptyData().Insert(); err != nil {
            err = gerror.Wrap(err, "新增失败，请稍后重试！")
        }
        return
    })
}
```

<!-- PATCHED: 2026-05-09 -->
**Every logic file MUST follow these strict constraints:**

**Header pattern** — every logic `.go` file must contain:
```go
type s{模块名} struct{}

func New{模块名}() *s{模块名} {
    return &s{模块名}{}
}

func init() {
    service.Register{模块名}(New{模块名}())
}
```
The struct name, constructor, and init function MUST use the same `{模块名}` (e.g. `sSysFoo`, `NewSysFoo`, `RegisterSysFoo`). These three blocks appear in this exact order at the top of the file after the `package` declaration and imports.

**Method signature pattern** — all methods MUST follow:
```go
func (s *s{模块名}) 方法名(ctx context.Context, in *input.{入参struct}) (out *input.{出参struct}, err error) {

}
```
- Receiver: `(s *s{模块名})` — always pointer receiver on the logic struct
- First parameter: `ctx context.Context` — always present
- Second parameter: `in *input.{入参struct}` — input DTO from `model/input/` (named `in`)
- Returns: `(out *input.{出参struct}, err error)` — named return values `out` and `err`; if no output DTO, omit `out`; if no input DTO (e.g. `Model` helper), omit `in`

<!-- PATCHED: 2026-05-09 -->
**All logic functions MUST be methods on the logic struct.** Standalone functions (`func functionName(...)`) are FORBIDDEN in logic files:

```go
// WRONG — standalone function
func doSomething(ctx context.Context) error {
    ...
}

// CORRECT — method on logic struct
func (s *sSysFoo) doSomething(ctx context.Context) (err error) {
    ...
}
```
The only allowed standalone functions in logic files are `New{模块名}()` and `init()` (the header pattern).

### 1.3 Logic Import Aggregator

TRIGGER: logic import, init aggregator, blank import

Location: `internal/logic/logic.go`

```go
package logic

import (
    _ "hotgo/internal/logic/admin"
    _ "hotgo/internal/logic/api"
    _ "hotgo/internal/logic/common"
    _ "hotgo/internal/logic/hook"
    _ "hotgo/internal/logic/middleware"
    _ "hotgo/internal/logic/sys"
    _ "hotgo/internal/logic/view"
)
```

> 输入 DTO 结构（`FooListInp`、`FooEditInp` 等）定义见 `gfstack-data` §3。DAO 模式和 ORM 操作方法（`dao.Foo`、`Page`、`ScanAndCount`、`WhereLike` 等）见 `gfstack-data` §4~§5。

---
name: gfstack-api
description: "HotGo API layer specification. Covers API request/response structs with g.Meta route binding, controller interface definitions, and controller implementation patterns. TRIGGER: API struct, g.Meta, route binding, request struct, response struct, path definition, Req/Res, controller interface, IModuleV1, controller implementation, thin handler, ControllerV1, NewV1. DO NOT TRIGGER: business logic (see gfstack-logic), data models (see gfstack-data)."
---

# gfstack-api

## 1. API Definition Layer

TRIGGER: API struct, g.Meta, route binding, request struct, response struct, path definition, Req/Res

Location: `api/{module}/v1/{entity}.go`

```go
package v1

import (
    "github.com/gogf/gf/v2/frame/g"
    "hotgo/internal/model/input/form"
    "hotgo/internal/model/input/sysin"
)

type FooListReq struct {
    g.Meta `path:"/foo/list" method:"get" tags:"Foo管理" summary:"获取Foo列表"`
    sysin.FooListInp
}

type FooListRes struct {
    form.PageRes
    List []*sysin.FooListModel `json:"list" dc:"数据列表"`
}

type FooEditReq struct {
    g.Meta `path:"/foo/edit" method:"post" tags:"Foo管理" summary:"修改/新增Foo"`
    sysin.FooEditInp
}

type FooEditRes struct{}

type FooViewReq struct {
    g.Meta `path:"/foo/view" method:"get" tags:"Foo管理" summary:"获取Foo详情"`
    sysin.FooViewInp
}

type FooViewRes struct {
    *sysin.FooViewModel
}
```

**Rules:**
- Req embeds an `*Inp` DTO from `model/input/`; Res embeds `form.PageRes` (lists) or `*ViewModel` (details)
- Collection field MUST be named `list`, never `items` or `data`
- Default sort: `Id DESC`
- JSON fields MUST include `dc` (description) tag
- `summary` and `tags` use Chinese labels
- Input DTO 结构定义见 `gfstack-data` §3；分页类型 `PageReq`/`PageRes` 见 `gfstack-data` §4

---

## 2. Controller Interface

TRIGGER: controller interface, IModuleV1, controller signature

Location: `api/{module}/{module}.go`

```go
package admin

import (
    "context"
    "hotgo/api/admin/v1"
)

type IAdminV1 interface {
    FooList(ctx context.Context, req *v1.FooListReq) (res *v1.FooListRes, err error)
    FooEdit(ctx context.Context, req *v1.FooEditReq) (res *v1.FooEditRes, err error)
    FooView(ctx context.Context, req *v1.FooViewReq) (res *v1.FooViewRes, err error)
}
```

---

## 3. Controller Implementation

TRIGGER: controller implementation, thin handler, ControllerV1, handler pattern, NewV1

Location: `internal/controller/{module}/{module}/{module}_v1_{entity}.go`

```go
package admin

import (
    "context"
    "hotgo/api/admin/v1"
    "hotgo/internal/model/input/sysin"
    "hotgo/internal/service"
)

type ControllerV1 struct{}

func NewV1() admin.IAdminV1 {
    return &ControllerV1{}
}

func (c *ControllerV1) FooList(ctx context.Context, req *v1.FooListReq) (res *v1.FooListRes, err error) {
    var (
        list       []*sysin.FooListModel
        totalCount int
    )
    if list, totalCount, err = service.SysFoo().List(ctx, &req.FooListInp); err != nil {
        return
    }
    if list == nil {
        list = []*sysin.FooListModel{}
    }
    res = new(v1.FooListRes)
    res.List = list
    res.PageRes.Pack(req, totalCount)
    return
}

func (c *ControllerV1) FooEdit(ctx context.Context, req *v1.FooEditReq) (res *v1.FooEditRes, err error) {
    err = service.SysFoo().Edit(ctx, &req.FooEditInp)
    return
}
```

<!-- PATCHED: 2026-05-09 -->
**CRITICAL: Controllers MUST be thin** — only call services and construct responses. NO business logic allowed. NO unused functions, variables, or constants — every declaration must be referenced.

> Service 接口定义和 Logic 实现模式见 `gfstack-logic`。控制器中调用的 `service.SysFoo().List()` / `service.SysFoo().Edit()` 的方法签名在 `gfstack-logic` §1 中定义。

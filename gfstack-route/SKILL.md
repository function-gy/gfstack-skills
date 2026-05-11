---
name: gfstack-route
description: "HotGo router and middleware specification. Covers route registration with group.Bind, middleware binding (AdminAuth, Develop), and middleware file conventions. TRIGGER: router, route registration, group.Bind, AdminAuth, route group, middleware binding, middleware, CORS, Ctx, ResponseHandler, Blacklist, PreFilter, RequestLog, DemoLimit, Develop. DO NOT TRIGGER: API definitions (see gfstack-api), business logic (see gfstack-logic)."
---

# gfstack-route

## 1. Router Registration

TRIGGER: router, route registration, group.Bind, AdminAuth, route group, middleware binding

Location: `internal/router/{module}.go`

```go
func Admin(ctx context.Context, group *ghttp.RouterGroup) {
    group.Group(simple.RouterPrefix(ctx, consts.AppAdmin), func(group *ghttp.RouterGroup) {
        group.Bind(
            admin.NewV1().SiteLogin,
            admin.NewV1().SitePing,
        )

        group.Middleware(service.Middleware().AdminAuth)
        group.Bind(
            admin.NewV1().FooList,
            admin.NewV1().FooView,
            admin.NewV1().FooEdit,
            admin.NewV1().FooDelete,
        )

        group.Middleware(service.Middleware().Develop)
        group.Bind(
            admin.NewV1().GenCodesList,
        )
    })
}
```

## 2. Middleware

TRIGGER: middleware, AdminAuth, CORS, Ctx, ResponseHandler, Blacklist, PreFilter, RequestLog, DemoLimit, Develop

All middleware in `internal/logic/middleware/`, registered via `service.RegisterMiddleware()`.

### Global Middleware Order (internal/cmd/http.go)

```go
s.BindMiddleware("/*any", []ghttp.HandlerFunc{
    service.Middleware().RequestLog,
    service.Middleware().Ctx,
    service.Middleware().CORS,
    service.Middleware().Blacklist,
    service.Middleware().DemoLimit,
    service.Middleware().ResponseHandler,
}...)
```

<!-- PATCHED: 2026-05-09 -->
**Middleware files MUST contain ONLY the main handler function.** Each middleware `.go` file should have exactly one function matching the pattern:

```go
func 中间件名称(r *ghttp.Request) {
    ...
}
```

- NO helper/auxiliary functions inside middleware files — all utilities, validators, data lookups belong in `internal/logic/{domain}/`
- Middleware is thin glue: receives request, calls logic services, decides to proceed (`r.Middleware.Next()`) or abort

> 请求全链路见 `gfstack-overview` §2。Controller 实现见 `gfstack-api` §3。

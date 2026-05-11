---
name: gfstack-infra
description: "Infrastructure specification. Covers token/JWT authentication system, HTTP server bootstrap (main.go), and standard JSON response format. TRIGGER: token, JWT, login, logout, ParseLoginUser, HS256, authentication, Bearer token, main.go, HTTP server, gcmd.Command, global.Init, server startup, Response struct, JSON response, response format, code message data. DO NOT TRIGGER: middleware (see gfstack-route), database (see gfstack-data)."
---

# gfstack-infra

## 1. Token System

TRIGGER: token, JWT, login, logout, ParseLoginUser, HS256, authentication, Bearer token

Location: `internal/library/token/token.go`

- Login: HS256 signed JWT + token metadata stored in cache
- Verification: parse JWT → retrieve cache metadata → check expiry → single-device detection → auto-refresh
- Cache keys: `token:{app}:{MD5(jwt)}` and `tokenBind:{app}:{userId}`
- config.yaml: secretKey, expires, autoRefresh, refreshInterval, maxRefreshTimes, multiLogin

## 2. HTTP Server Bootstrap

TRIGGER: main.go, HTTP server, gcmd.Command, global.Init, server startup

```go
func main() {
    var ctx = gctx.GetInitCtx()
    global.Init(ctx)
    cmd.Main.Run(ctx)
}
```

## 3. Response Format

TRIGGER: Response struct, JSON response, response format, code message data

```go
type Response struct {
    Code      int         `json:"code"      dc:"状态码"`
    Message   string      `json:"message,omitempty" dc:"提示信息"`
    Data      interface{} `json:"data,omitempty"    dc:"响应数据"`
    Error     interface{} `json:"error,omitempty"   dc:"错误详情"`
    Timestamp int64       `json:"timestamp" dc:"响应时间戳"`
    TraceID   string      `json:"traceID"   dc:"链路追踪ID"`
}
```

> Token 校验中间件 AdminAuth 见 `gfstack-route` §2。错误码定义和错误处理规范见 `gfstack-style` §1。

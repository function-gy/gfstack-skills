// ================================================================
// 示例: Contexts 请求上下文 (internal/library/contexts/context.go)
// 在 HTTP 请求生命周期中存储和获取用户信息、响应数据、模块名等
// ================================================================

package contexts

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/internal/consts"
	"hotgo/internal/model"
)

// Init 初始化上下文对象指针到 ghttp.Request 上下文
func Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(consts.ContextHTTPKey, customCtx)
}

// Get 获得上下文变量，没有设置返回 nil
func Get(ctx context.Context) *model.Context {
	value := ctx.Value(consts.ContextHTTPKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

// SetUser 将用户信息设置到上下文
func SetUser(ctx context.Context, user *model.Identity) {
	c := Get(ctx)
	if c == nil {
		return
	}
	c.User = user
}

// GetUser 获取用户信息
func GetUser(ctx context.Context) *model.Identity {
	c := Get(ctx)
	if c == nil {
		return nil
	}
	return c.User
}

// GetUserId 获取用户ID
func GetUserId(ctx context.Context) int64 {
	user := GetUser(ctx)
	if user == nil {
		return 0
	}
	return user.Id
}

// GetRoleId 获取用户角色ID
func GetRoleId(ctx context.Context) int64 {
	user := GetUser(ctx)
	if user == nil {
		return 0
	}
	return user.RoleId
}

// GetRoleKey 获取用户角色唯一编码
func GetRoleKey(ctx context.Context) string {
	user := GetUser(ctx)
	if user == nil {
		return ""
	}
	return user.RoleKey
}

// GetModule 获取当前请求的应用模块（admin|api|default）
func GetModule(ctx context.Context) string {
	c := Get(ctx)
	if c == nil {
		return ""
	}
	return c.Module
}

// SetModule 设置应用模块
func SetModule(ctx context.Context, module string) {
	c := Get(ctx)
	if c == nil {
		return
	}
	c.Module = module
}

// SetResponse 设置响应数据，用于访问日志
func SetResponse(ctx context.Context, response *model.Response) {
	c := Get(ctx)
	if c == nil {
		return
	}
	c.Response = response
}

// SetData 设置额外 KV 数据
func SetData(ctx context.Context, k string, v interface{}) {
	c := Get(ctx)
	if c == nil {
		return
	}
	c.Data[k] = v
}

// GetData 获取额外数据
func GetData(ctx context.Context) g.Map {
	c := Get(ctx)
	if c == nil {
		return nil
	}
	return c.Data
}

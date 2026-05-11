// ================================================================
// 示例8: 路由注册 (internal/router/admin.go)
// 路径: internal/router/{module}.go
// 作用: 绑定 Controller 方法到路由组，区分公开/认证/开发路由
// ================================================================

package router

import (
	"context"

	"hotgo/internal/consts"
	"hotgo/internal/controller/admin/admin"
	"hotgo/internal/service"
	"hotgo/utility/simple"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Admin(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, consts.AppAdmin), func(group *ghttp.RouterGroup) {

		// ===== 第一组: 公开路由（无需认证） =====
		group.Bind(
			admin.NewV1().SiteLoginLogout,  // 登录/退出
			admin.NewV1().SiteRegister,     // 注册
			admin.NewV1().SiteLoginCaptcha, // 验证码
			admin.NewV1().SiteAccountLogin, // 账号登录
			admin.NewV1().SitePing,         // 健康检查
		)

		// ===== 第二组: 需要认证的路由 =====
		group.Middleware(service.Middleware().AdminAuth)
		group.Bind(
			admin.NewV1().FooList,
			admin.NewV1().FooView,
			admin.NewV1().FooEdit,
			admin.NewV1().FooDelete,
			admin.NewV1().FooSwitch,
			// ... 其他需要认证的路由
		)

		// ===== 第三组: 开发工具路由（IP白名单过滤） =====
		group.Middleware(service.Middleware().Develop)
		group.Bind(
			admin.NewV1().GenCodesList,
			admin.NewV1().GenCodesView,
			admin.NewV1().GenCodesEdit,
			admin.NewV1().GenCodesDelete,
			admin.NewV1().GenCodesBuild,
		)
	})
}

// 关键规则:
// 1. group.Bind() 一次性绑定多个方法，每个方法签名必须为 (ctx, *Req) (*Res, error)
// 2. 分组顺序: 公开路由 → 认证路由 → 开发工具路由
// 3. Middleware 在 Bind 之前设置，才能对该组内路由生效

// ================================================================
// 示例: Token 登录与验证 (internal/library/token/token.go)
// 代码来自 HS_COUPON 项目，展示完整的 JWT token 生成、验证、刷新机制
// ================================================================

package token

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"

	"hotgo/internal/consts"
	"hotgo/internal/library/cache"
	"hotgo/internal/library/contexts"
	"hotgo/internal/model"
	"hotgo/utility/simple"
)

// ================================================================
// Claims — JWT 载荷，嵌入 Identity 用户身份
// ================================================================

type Claims struct {
	*model.Identity
	jwt.RegisteredClaims
}

// ================================================================
// Token — 缓存的 token 元数据
// ================================================================

type Token struct {
	ExpireAt     int64 `json:"exp"` // token过期时间
	RefreshAt    int64 `json:"ra"`  // 刷新时间
	RefreshCount int64 `json:"rc"`  // 刷新次数
}

// ================================================================
// TokenConfig — token 配置结构体（对应 config.yaml 中 token 段）
// ================================================================
//
// config.yaml 配置示例:
//
//	token:
//	  secretKey: "w37DDEvspVx4Bg5t"   # 令牌加密秘钥
//	  expires: 300                    # 令牌有效期，单位：秒
//	  autoRefresh: true               # 是否开启自动刷新过期时间
//	  refreshInterval: 150            # 刷新间隔，单位：秒
//	  maxRefreshTimes: -1             # 最大允许刷新次数，-1不限制
//	  multiLogin: true                # 是否允许多端登录

type TokenConfig struct {
	SecretKey       string `json:"secretKey"`
	Expires         int64  `json:"expires"`
	AutoRefresh     bool   `json:"autoRefresh"`
	RefreshInterval int64  `json:"refreshInterval"`
	MaxRefreshTimes int64  `json:"maxRefreshTimes"`
	MultiLogin      bool   `json:"multiLogin"`
}

// ================================================================
// 全局变量
// ================================================================

var (
	config          *TokenConfig
	errorLogin      = gerror.New("登录身份已失效，请重新登录！")
	errorMultiLogin = gerror.New("账号已在其他地方登录，如非本人操作请及时修改登录密码！")
)

func SetConfig(c *TokenConfig) {
	config = c
}

func GetConfig() *TokenConfig {
	return config
}

// ================================================================
// Login — 登录：生成 JWT token 并缓存
// ================================================================

func Login(ctx context.Context, user *model.Identity) (string, int64, error) {
	claims := Claims{
		user,
		jwt.RegisteredClaims{},
	}

	// 1. 使用 HS256 生成 JWT
	header, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", 0, err
	}

	var (
		now     = gtime.Now()
		authKey = GetAuthKey(header)                        // 认证key（token的MD5）
		tokenKey = GetTokenKey(user.App, authKey)            // 缓存key: token:{app}:{authKey}
		bindKey  = GetBindKey(user.App, user.Id)             // 绑定key: tokenBind:{app}:{userId}
		duration = time.Second * gconv.Duration(config.Expires)
	)

	token := &Token{
		ExpireAt:     now.Unix() + config.Expires,
		RefreshAt:    now.Unix(),
		RefreshCount: 0,
	}

	// 2. 缓存 token
	if err = cache.Instance().Set(ctx, tokenKey, token, duration); err != nil {
		return "", 0, err
	}

	// 3. 缓存用户绑定（用于单端登录检测）
	if err = cache.Instance().Set(ctx, bindKey, tokenKey, duration); err != nil {
		return "", 0, err
	}

	return header, config.Expires, nil
}

// ================================================================
// Logout — 注销登录：删除缓存中的 token
// ================================================================

func Logout(r *ghttp.Request) (err error) {
	var (
		ctx    = r.Context()
		header = GetAuthorization(r)
	)

	if header == "" {
		return errorLogin
	}

	claims, err := parseToken(ctx, header)
	if err != nil {
		return errorLogin
	}

	var (
		authKey  = GetAuthKey(header)
		tokenKey = GetTokenKey(contexts.GetModule(ctx), authKey)
		bindKey  = GetBindKey(contexts.GetModule(ctx), claims.Id)
	)

	// 删除 token
	if _, err = cache.Instance().Remove(ctx, tokenKey); err != nil {
		return
	}

	// 如果不是多端登录，也删除绑定
	if !config.MultiLogin {
		if _, err = cache.Instance().Remove(ctx, bindKey); err != nil {
			return
		}
	}
	return
}

// ================================================================
// ParseLoginUser — 解析登录用户信息（中间件调用入口）
// ================================================================

func ParseLoginUser(r *ghttp.Request) (user *model.Identity, err error) {
	var (
		ctx    = r.Context()
		header = GetAuthorization(r)
	)

	if header == "" {
		return nil, errorLogin
	}

	// 1. 解析 JWT token
	claims, err := parseToken(ctx, header)
	if err != nil {
		return nil, errorLogin
	}

	var (
		authKey  = GetAuthKey(header)
		tokenKey = GetTokenKey(claims.App, authKey)
		bindKey  = GetBindKey(claims.App, claims.Id)
	)

	// 2. 从缓存获取 token 元数据
	tk, err := cache.Instance().Get(ctx, tokenKey)
	if err != nil || tk.IsEmpty() {
		return nil, errorLogin
	}

	var token *Token
	if err = tk.Scan(&token); err != nil || token == nil {
		return nil, errorLogin
	}

	// 3. 检查是否过期
	now := gtime.Now()
	if token.ExpireAt < now.Unix() {
		return nil, errorLogin
	}

	// 4. 单端登录检测
	if !config.MultiLogin {
		origin, err := cache.Instance().Get(ctx, bindKey)
		if err != nil || origin == nil || origin.IsEmpty() {
			return nil, errorLogin
		}
		if tokenKey != origin.String() {
			return nil, errorMultiLogin
		}
	}

	// 5. 自动刷新 token 有效期
	simple.SafeGo(ctx, func(ctx context.Context) {
		if !config.AutoRefresh {
			return
		}
		if config.MaxRefreshTimes != -1 && token.RefreshCount >= config.MaxRefreshTimes {
			return
		}
		if gtime.New(token.RefreshAt).Unix()+config.RefreshInterval > now.Unix() {
			return
		}

		token.ExpireAt = now.Unix() + config.Expires
		token.RefreshAt = now.Unix()
		token.RefreshCount += 1
		duration := time.Second * gconv.Duration(config.Expires)

		cache.Instance().Set(ctx, tokenKey, token, duration)
		cache.Instance().Set(ctx, bindKey, tokenKey, duration)
	})

	return claims.Identity, nil
}

// ================================================================
// 内部辅助函数
// ================================================================

// parseToken 解析 JWT 令牌
func parseToken(ctx context.Context, header string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(header, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errorLogin
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errorLogin
	}
	return claims, nil
}

// GetAuthorization 从请求头或URL参数获取 Authorization
func GetAuthorization(r *ghttp.Request) string {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return r.Get("authorization").String()
	}
	return gstr.Replace(authorization, "Bearer ", "")
}

// GetAuthKey 认证key（对 token 做 MD5）
func GetAuthKey(token string) string {
	return gmd5.MustEncryptString("hotgo" + token)
}

// GetTokenKey 令牌缓存key 格式: token:{app}:{authKey}
func GetTokenKey(appName, authKey string) string {
	return fmt.Sprintf("%v:%v:%v", consts.CacheToken, appName, authKey)
}

// GetBindKey 令牌身份绑定key 格式: tokenBind:{app}:{userId}
func GetBindKey(appName string, userId int64) string {
	return fmt.Sprintf("%v:%v:%v", consts.CacheTokenBind, appName, userId)
}

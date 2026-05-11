// ================================================================
// 示例: 中间件全集 (internal/logic/middleware/)
// 包含所有必需的中间件，代码来自 HS_COUPON 项目
// ================================================================

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"

	"hotgo/internal/consts"
	"hotgo/internal/global"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/location"
	"hotgo/internal/library/response"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	"hotgo/internal/service"
	"hotgo/utility/charset"
	"hotgo/utility/simple"
	"hotgo/utility/validate"
)

// ================================================================
// sMiddleware 中间件结构体
// ================================================================

type sMiddleware struct {
	LoginUrl         string // 登录路由地址
	DemoWhiteList    g.Map  // 演示模式放行的路由白名单
	NotRecordRequest g.Map  // 不记录请求数据的路由
}

func NewMiddleware() *sMiddleware {
	return &sMiddleware{
		LoginUrl: "/common",
		DemoWhiteList: g.Map{
			"/admin/site/accountLogin": struct{}{},
			"/admin/site/mobileLogin":  struct{}{},
		},
		NotRecordRequest: g.Map{
			"/admin/upload/file":       struct{}{},
			"/admin/upload/uploadPart": struct{}{},
		},
	}
}

func init() {
	service.RegisterMiddleware(NewMiddleware())
}

// ================================================================
// 1. Ctx — 请求上下文初始化（必须第一个执行）
// ================================================================

func (s *sMiddleware) Ctx(r *ghttp.Request) {
	// 国际化
	r.SetCtx(gi18n.WithLanguage(r.Context(), simple.GetHeaderLocale(r.Context())))

	data := make(g.Map)
	if _, ok := s.NotRecordRequest[r.URL.Path]; ok {
		data["request.body"] = gjson.New(nil)
	} else {
		data["request.body"] = gjson.New(r.GetBodyString())
	}

	contexts.Init(r, &model.Context{
		Data:   data,
		Module: getModule(r.URL.Path),
	})

	if len(r.Cookie.GetSessionId()) == 0 {
		r.Cookie.SetSessionId(gctx.CtxId(r.Context()))
	}

	r.SetCtx(r.GetNeverDoneCtx())
	r.Middleware.Next()
}

func getModule(path string) (module string) {
	slice := strings.Split(path, "/")
	if len(slice) < 2 || slice[1] == "" {
		return consts.AppDefault
	}
	return slice[1]
}

// ================================================================
// 2. CORS — 跨域中间件
// ================================================================

func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// ================================================================
// 3. Blacklist — IP黑名单限制中间件
// ================================================================

func (s *sMiddleware) Blacklist(r *ghttp.Request) {
	if err := service.SysBlacklist().VerifyRequest(r); err != nil {
		response.JsonExit(r, gerror.Code(err).Code(), err.Error())
	}
	r.Middleware.Next()
}

// ================================================================
// 4. DemoLimit — 演示系统操作限制（仅当 isDemo=true 时生效）
// ================================================================

func (s *sMiddleware) DemoLimit(r *ghttp.Request) {
	if !simple.IsDemo(r.Context()) {
		r.Middleware.Next()
		return
	}
	if r.Method == http.MethodPost {
		if _, ok := s.DemoWhiteList[r.URL.Path]; ok {
			r.Middleware.Next()
			return
		}
		response.JsonExit(r, gcode.CodeNotSupported.Code(), "演示系统禁止操作！")
		return
	}
	r.Middleware.Next()
}

// ================================================================
// 5. Develop — 开发工具白名单过滤（IP白名单）
// ================================================================

func (s *sMiddleware) Develop(r *ghttp.Request) {
	ips := g.Cfg().MustGet(r.Context(), "hggen.allowedIPs").Strings()
	if len(ips) == 0 {
		response.JsonExit(r, gcode.CodeNotSupported.Code(), "请配置生成白名单！")
		return
	}
	if !gstr.InArray(ips, "*") {
		clientIp := location.GetClientIp(r)
		ok := false
		for _, ip := range ips {
			if ip == clientIp {
				ok = true
				break
			}
		}
		if !ok {
			response.JsonExit(r, gcode.CodeNotSupported.Code(),
				fmt.Sprintf("当前IP[%s]没有配置生成白名单！", clientIp))
			return
		}
	}
	r.Middleware.Next()
}

// ================================================================
// 6. PreFilter — 请求输入预处理（自动调用 Filter() 验证）
// ================================================================

func (s *sMiddleware) PreFilter(r *ghttp.Request) {
	router := global.GetRequestRoute(r)
	if router == nil {
		r.Middleware.Next()
		return
	}

	funcInfo := router.Handler.Info
	if funcInfo.Type.NumIn() != 2 {
		r.Middleware.Next()
		return
	}

	inputType := funcInfo.Type.In(1)
	var inputObject reflect.Value
	if inputType.Kind() == reflect.Ptr {
		inputObject = reflect.New(inputType.Elem())
	} else {
		inputObject = reflect.New(inputType.Elem()).Elem()
	}

	if err := r.Parse(inputObject.Interface()); err != nil {
		resp := gerror.Code(err)
		response.JsonExit(r, resp.Code(), gerror.Current(err).Error(), resp.Detail())
		return
	}

	if _, ok := inputObject.Interface().(validate.Filter); !ok {
		r.Middleware.Next()
		return
	}

	if err := validate.PreFilter(r.Context(), inputObject.Interface()); err != nil {
		resp := gerror.Code(err)
		response.JsonExit(r, resp.Code(), gerror.Current(err).Error(), resp.Detail())
		return
	}

	r.SetParamMap(gconv.Map(inputObject.Interface()))
	r.Middleware.Next()
}

// ================================================================
// 7. ResponseHandler — HTTP响应统一处理
// ================================================================

func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	switch r.Response.Status {
	case 403:
		r.Response.Writeln("403 - 网站拒绝显示此网页")
		return
	case 404:
		r.Response.Writeln("404 - 你似乎来到了没有知识存在的荒原…")
		return
	}

	contentType := getContentType(r)
	if contentType != consts.HTTPContentTypeStream && r.Response.BufferLength() > 0 {
		return
	}

	switch contentType {
	case consts.HTTPContentTypeHtml:
		s.responseHtml(r)
	case consts.HTTPContentTypeXml:
		s.responseXml(r)
	default:
		responseJson(r)
	}
}

func (s *sMiddleware) responseHtml(r *ghttp.Request) {
	code, message, resp := parseResponse(r)
	if code == gcode.CodeOK.Code() {
		return
	}
	r.Response.ClearBuffer()
	_ = r.Response.WriteTplContent(simple.DefaultErrorTplContent(r.Context()),
		g.Map{"code": code, "message": message, "stack": resp})
}

func (s *sMiddleware) responseXml(r *ghttp.Request) {
	code, message, data := parseResponse(r)
	response.RXml(r, code, message, data)
}

func responseJson(r *ghttp.Request) {
	code, message, data := parseResponse(r)
	response.RJson(r, code, message, data)
}

func parseResponse(r *ghttp.Request) (code int, message string, resp interface{}) {
	ctx := r.Context()
	err := r.GetError()
	if err == nil {
		return gcode.CodeOK.Code(), "操作成功", r.GetHandlerResponse()
	}

	if simple.Debug(ctx) {
		message = gerror.Current(err).Error()
		if getContentType(r) == consts.HTTPContentTypeHtml {
			resp = charset.SerializeStack(err)
		} else {
			resp = charset.ParseErrStack(err)
		}
	} else {
		message = consts.ErrorMessage(gerror.Current(err))
	}

	code = gerror.Code(err).Code()
	if code == gcode.CodeNil.Code() {
		g.Log().Stdout(false).Infof(ctx, "exception:%v", err)
	} else {
		g.Log().Errorf(ctx, "exception:%v", err)
	}
	return
}

func getContentType(r *ghttp.Request) (contentType string) {
	contentType = r.Response.Header().Get("Content-Type")
	if contentType != "" {
		return
	}
	mime := gmeta.Get(r.GetHandlerResponse(), "mime").String()
	if mime == "" {
		contentType = consts.HTTPContentTypeJson
	} else {
		contentType = mime
	}
	return
}

// ================================================================
// 8. AdminAuth — 后台鉴权中间件
// ================================================================

func (s *sMiddleware) AdminAuth(r *ghttp.Request) {
	var (
		ctx  = r.Context()
		path = gstr.Replace(r.URL.Path, simple.RouterPrefix(ctx, consts.AppAdmin), "", 1)
	)

	// 不需要验证登录的路由地址
	if s.IsExceptLogin(ctx, consts.AppAdmin, path) {
		r.Middleware.Next()
		return
	}

	// 将用户信息传递到上下文中
	if err := s.DeliverUserContext(r); err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), err.Error())
		return
	}

	// 不需要验证权限的路由地址
	if s.IsExceptAuth(ctx, consts.AppAdmin, path) {
		r.Middleware.Next()
		return
	}

	r.Middleware.Next()
}

// ================================================================
// 9. ApiSign — 对外API签名验证中间件
// ================================================================

func (s *sMiddleware) ApiSign(r *ghttp.Request) {
	var (
		ctx        = r.Context()
		appId      = r.Header.Get("X-Appid")
		xSignature = r.Header.Get("X-Signature")
		params     = make(map[string]interface{})
	)

	if !g.IsEmpty(appId) && appId == "testAppid" {
		goto Next
	}
	if g.IsEmpty(appId) || g.IsEmpty(xSignature) {
		response.JsonExit(r, gcode.CodeInvalidRequest.Code(), "缺少Auth参数")
		return
	}
	if appId != g.Cfg().MustGet(ctx, "sign.signKey").String() {
		response.JsonExit(r, gcode.CodeSecurityReason.Code(), "你没有访问权限！")
		return
	}

	// 根据请求方法处理参数
	switch r.Method {
	case "GET":
		for k, v := range r.GetQueryMap() {
			params[k] = v
		}
	case "POST", "PUT", "PATCH":
		contentType := r.Header.Get("Content-Type")
		if gstr.Contains(contentType, "multipart/form-data") {
			for k, v := range r.GetMultipartForm().Value {
				params[k] = v[0]
			}
		} else {
			jsonData, _ := r.GetJson()
			if jsonData != nil {
				for k, v := range gconv.Map(jsonData) {
					if g.IsNil(v) || reflect.TypeOf(v).Kind() == reflect.Bool {
						continue
					}
					params[k] = v
				}
			}
		}
	}

	if xSignature != SignStr(params, g.Cfg().MustGet(ctx, "sign.signSecret").String(), "key", "X-Signature") {
		response.JsonExit(r, gcode.CodeInvalidRequest.Code(), "权限验证失败")
	}
Next:
	r.Middleware.Next()
}

// ================================================================
// 辅助方法
// ================================================================

// DeliverUserContext 将用户信息传递到上下文中
func (s *sMiddleware) DeliverUserContext(r *ghttp.Request) (err error) {
	user, err := token.ParseLoginUser(r)
	if err != nil {
		return
	}
	switch user.App {
	case consts.AppAdmin:
		if err = service.AdminSite().BindUserContext(r.Context(), user); err != nil {
			return
		}
	default:
		contexts.SetUser(r.Context(), user)
	}
	return
}

// IsExceptAuth 是否是不需要验证权限的路由地址
func (s *sMiddleware) IsExceptAuth(ctx context.Context, appName, path string) bool {
	pathList := g.Cfg().MustGet(ctx, fmt.Sprintf("router.%v.exceptAuth", appName)).Strings()
	for i := 0; i < len(pathList); i++ {
		if validate.InSliceExistStr(pathList[i], path) {
			return true
		}
	}
	return false
}

// IsExceptLogin 是否是不需要登录的路由地址
func (s *sMiddleware) IsExceptLogin(ctx context.Context, appName, path string) bool {
	pathList := g.Cfg().MustGet(ctx, fmt.Sprintf("router.%v.exceptLogin", appName)).Strings()
	for i := 0; i < len(pathList); i++ {
		if validate.InSliceExistStr(pathList[i], path) {
			return true
		}
	}
	return false
}

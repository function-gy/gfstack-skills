// ================================================================
// 示例: Response 响应工具 (internal/library/response/response.go)
// 统一的 JSON/XML 响应封装，标准格式: {code, message, data, error, timestamp, traceID}
// ================================================================

package response

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/internal/library/contexts"
	"hotgo/internal/model"
)

// JsonExit 返回JSON数据并退出当前 HTTP 执行函数
func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	RJson(r, code, message, data...)
	r.Exit()
}

// RJson 标准 JSON 响应
func RJson(r *ghttp.Request, code int, message string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	res := &model.Response{
		Code:      code,
		Message:   message,
		Timestamp: gtime.Timestamp(),
		TraceID:   gctx.CtxId(r.Context()),
	}
	if gcode.CodeOK.Code() == code {
		res.Data = responseData
	} else {
		res.Error = responseData
	}

	r.Response.ClearBuffer()
	r.Response.WriteJson(res)
	contexts.SetResponse(r.Context(), res)
}

// RXml 标准 XML 响应
func RXml(r *ghttp.Request, code int, message string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	res := &model.Response{
		Code:      code,
		Message:   message,
		Timestamp: gtime.Timestamp(),
		TraceID:   gctx.CtxId(r.Context()),
	}
	if gcode.CodeOK.Code() == code {
		res.Data = responseData
	} else {
		res.Error = responseData
	}

	r.Response.ClearBuffer()
	r.Response.WriteXml(gconv.Map(res))
	contexts.SetResponse(r.Context(), res)
}

// CustomJson 自定义 JSON（直接透传内容）
func CustomJson(r *ghttp.Request, content interface{}) {
	r.Response.ClearBuffer()
	r.Response.WriteJson(content)
	contexts.SetResponse(r.Context(), &model.Response{
		Code:      0,
		Data:      content,
		Timestamp: gtime.Timestamp(),
		TraceID:   gctx.CtxId(r.Context()),
	})
}

// Redirect 重定向
func Redirect(r *ghttp.Request, location string, code ...int) {
	r.Response.RedirectTo(location, code...)
}

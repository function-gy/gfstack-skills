// ================================================================
// 示例1: API 定义层 (api/admin/v1/foo.go)
// 路径: api/{module}/v1/{entity}.go
// 作用: 定义 HTTP 请求/响应结构体，使用 g.Meta 自动注册路由
// ================================================================

package v1

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

// ---------- 列表 ----------

type FooListReq struct {
	g.Meta `path:"/foo/list" method:"get" tags:"Foo管理" summary:"获取Foo列表"`
	sysin.FooListInp
}

type FooListRes struct {
	form.PageRes
	List []*sysin.FooListModel `json:"list" dc:"数据列表"`
}

// ---------- 查看 ----------

type FooViewReq struct {
	g.Meta `path:"/foo/view" method:"get" tags:"Foo管理" summary:"获取Foo详情"`
	sysin.FooViewInp
}

type FooViewRes struct {
	*sysin.FooViewModel
}

// ---------- 编辑（新增/修改） ----------

type FooEditReq struct {
	g.Meta `path:"/foo/edit" method:"post" tags:"Foo管理" summary:"修改/新增Foo"`
	sysin.FooEditInp
}

type FooEditRes struct{}

// ---------- 删除 ----------

type FooDeleteReq struct {
	g.Meta `path:"/foo/delete" method:"post" tags:"Foo管理" summary:"删除Foo"`
	sysin.FooDeleteInp
}

type FooDeleteRes struct{}

// ---------- 状态开关 ----------

type FooSwitchReq struct {
	g.Meta `path:"/foo/switch" method:"post" tags:"Foo管理" summary:"更新Foo状态"`
	sysin.FooSwitchInp
}

type FooSwitchRes struct{}

// ================================================================
// 示例5: Input DTO + 字段过滤 (internal/model/input/sysin/foo.go)
// 路径: internal/model/input/{domain}in/{entity}.go
// 作用: 定义请求入参、响应模型、字段白名单过滤、Filter 验证
// ================================================================

package sysin

import (
	"context"

	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ======================== 字段过滤结构体 ========================

// FooUpdateFields 修改时允许更新的字段（白名单）
type FooUpdateFields struct {
	Name   string `json:"name"   dc:"名称"`
	Status string `json:"status" dc:"状态"`
	Remark string `json:"remark" dc:"备注"`
}

// FooInsertFields 新增时允许插入的字段（白名单）
type FooInsertFields struct {
	Name   string `json:"name"   dc:"名称"`
	Status string `json:"status" dc:"状态"`
	Remark string `json:"remark" dc:"备注"`
}

// ======================== 列表 ========================

// FooListInp 列表查询入参
type FooListInp struct {
	form.PageReq
	Name   string `json:"name"   dc:"名称"`
	Status string `json:"status" dc:"状态"`
}

// Filter 自定义验证规则（列表通常无需验证，留空即可）
func (in *FooListInp) Filter(ctx context.Context) (err error) {
	return
}

// FooListModel 列表返回项
type FooListModel struct {
	Id        int64       `json:"id"        dc:"id"`
	Name      string      `json:"name"      dc:"名称"`
	Status    string      `json:"status"    dc:"状态"`
	Remark    string      `json:"remark"    dc:"备注"`
	CreatedAt *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// ======================== 编辑（新增/修改二合一） ========================

// FooEditInp 编辑入参 — 嵌入 entity 继承所有 DB 字段
type FooEditInp struct {
	entity.Foo
}

// Filter 编辑时的自定义验证规则
func (in *FooEditInp) Filter(ctx context.Context) (err error) {
	// 验证名称
	if err := g.Validator().Rules("required").Data(in.Name).Messages("名称不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	// 验证状态
	if err := g.Validator().Rules("required").Data(in.Status).Messages("状态不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	return
}

// ======================== 查看 ========================

// FooViewInp 查看详情入参
type FooViewInp struct {
	Id int64 `json:"id" v:"required#id不能为空" dc:"id"`
}

func (in *FooViewInp) Filter(ctx context.Context) (err error) {
	return
}

// FooViewModel 查看详情返回 — 嵌入 entity 获取所有字段
type FooViewModel struct {
	entity.Foo
}

// ======================== 删除 ========================

// FooDeleteInp 删除入参
type FooDeleteInp struct {
	Id interface{} `json:"id" v:"required#id不能为空" dc:"id"`
}

// ======================== 状态开关 ========================

// FooSwitchInp 状态开关入参
type FooSwitchInp struct {
	form.SwitchReq // 嵌入通用开关结构体（含 Status 字段）
	Id int64       `json:"id" v:"required#id不能为空" dc:"id"`
}

func (in *FooSwitchInp) Filter(ctx context.Context) (err error) {
	return
}

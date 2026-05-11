// ================================================================
// 示例2: Controller 实现 (internal/controller/admin/admin/admin_v1_foo.go)
// 路径: internal/controller/{module}/{module}/{module}_v1_{entity}.go
// 作用: 薄 Handler，只调用 Service，构造响应
// ================================================================

package admin

import (
	"context"

	"hotgo/api/admin/v1"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

type ControllerV1 struct{}

// NewV1 返回 Controller 接口实例
func NewV1() admin.IAdminV1 { // interface 定义在 api/admin/admin.go
	return &ControllerV1{}
}

// FooList 列表
func (c *ControllerV1) FooList(ctx context.Context, req *v1.FooListReq) (res *v1.FooListRes, err error) {
	list, totalCount, err := service.SysFoo().List(ctx, &req.FooListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.FooListModel{} // 关键: 永远不要返回 nil
	}
	res = new(v1.FooListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount) // 填充分页信息
	return
}

// FooView 查看
func (c *ControllerV1) FooView(ctx context.Context, req *v1.FooViewReq) (res *v1.FooViewRes, err error) {
	data, err := service.SysFoo().View(ctx, &req.FooViewInp)
	if err != nil {
		return
	}
	res = new(v1.FooViewRes)
	res.FooViewModel = data
	return
}

// FooEdit 编辑（新增/修改二合一）
func (c *ControllerV1) FooEdit(ctx context.Context, req *v1.FooEditReq) (res *v1.FooEditRes, err error) {
	err = service.SysFoo().Edit(ctx, &req.FooEditInp)
	return
}

// FooDelete 删除
func (c *ControllerV1) FooDelete(ctx context.Context, req *v1.FooDeleteReq) (res *v1.FooDeleteRes, err error) {
	err = service.SysFoo().Delete(ctx, &req.FooDeleteInp)
	return
}

// FooSwitch 状态开关
func (c *ControllerV1) FooSwitch(ctx context.Context, req *v1.FooSwitchReq) (res *v1.FooSwitchRes, err error) {
	err = service.SysFoo().Switch(ctx, &req.FooSwitchInp)
	return
}

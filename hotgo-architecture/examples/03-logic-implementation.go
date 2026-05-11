// ================================================================
// 示例4: Logic 实现 (internal/logic/sys/foo.go)
// 路径: internal/logic/{domain}/{entity}.go
// 作用: 业务逻辑具体实现，通过 init() 自注册到 Service
// ================================================================

package sys

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/library/hgorm"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// sSysFoo 实现 ISysFoo 接口
type sSysFoo struct{}

func NewSysFoo() *sSysFoo {
	return &sSysFoo{}
}

// init 自动注册到 Service 层
func init() {
	service.RegisterSysFoo(NewSysFoo())
}

// Model 返回带权限过滤的 ORM 模型
func (s *sSysFoo) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.Foo.Ctx(ctx), option...)
}

// List 获取列表
func (s *sSysFoo) List(ctx context.Context, in *sysin.FooListInp) (list []*sysin.FooListModel, totalCount int, err error) {
	mod := s.Model(ctx)

	// 字段白名单过滤
	mod = mod.Fields(sysin.FooListModel{})

	// 条件筛选
	if in.Name != "" {
		mod = mod.WhereLike(dao.Foo.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status != "" {
		mod = mod.Where(dao.Foo.Columns().Status, in.Status)
	}

	// 分页 + 排序
	mod = mod.Page(in.Page, in.PerPage)
	mod = mod.OrderDesc(dao.Foo.Columns().Id)

	if err = mod.ScanAndCount(&list, &totalCount, false); err != nil {
		err = gerror.Wrap(err, "获取Foo列表失败，请稍后重试！")
		return
	}
	return
}

// Edit 修改/新增
func (s *sSysFoo) Edit(ctx context.Context, in *sysin.FooEditInp) (err error) {
	// 1. 参数验证（触发 Filter 方法）
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}

	// 2. 唯一性校验
	if err = hgorm.IsUnique(ctx, &dao.Foo, g.Map{dao.Foo.Columns().Name: in.Name}, "名称已存在", in.Id); err != nil {
		return
	}

	// 3. 事务操作
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		// 修改
		if in.Id > 0 {
			if _, err = s.Model(ctx).
				Fields(sysin.FooUpdateFields{}). // 只更新允许的字段
				WherePri(in.Id).
				Data(in).
				Update(); err != nil {
				err = gerror.Wrap(err, "修改Foo失败，请稍后重试！")
			}
			return
		}

		// 新增（需要 FilterAuth: false 因为新建时无用户权限上下文）
		if _, err = s.Model(ctx, &handler.Option{FilterAuth: false}).
			Fields(sysin.FooInsertFields{}). // 只插入允许的字段
			Data(in).
			OmitEmptyData(). // 跳过空值字段
			Insert(); err != nil {
			err = gerror.Wrap(err, "新增Foo失败，请稍后重试！")
		}
		return
	})
}

// View 获取详情
func (s *sSysFoo) View(ctx context.Context, in *sysin.FooViewInp) (res *sysin.FooViewModel, err error) {
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}
	res = new(sysin.FooViewModel)
	if err = s.Model(ctx).
		Fields(sysin.FooViewModel{}).
		WherePri(in.Id).
		Scan(res); err != nil {
		err = gerror.Wrap(err, "获取Foo详情失败，请稍后重试！")
		return
	}
	return
}

// Delete 删除（软删除，框架自动设置 deleted_at）
func (s *sSysFoo) Delete(ctx context.Context, in *sysin.FooDeleteInp) (err error) {
	// 注意: 不需要手动 validate.PreFilter，因为可以用 struct tag 验证
	if _, err = dao.Foo.Ctx(ctx).WherePri(in.Id).Delete(); err != nil {
		err = gerror.Wrap(err, "删除Foo失败，请稍后重试！")
	}
	return
}

// Switch 状态开关
func (s *sSysFoo) Switch(ctx context.Context, in *sysin.FooSwitchInp) (err error) {
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}
	if _, err = s.Model(ctx).
		Fields(dao.Foo.Columns().Status).
		WherePri(in.Id).
		Data(g.Map{dao.Foo.Columns().Status: in.Status}).
		Update(); err != nil {
		err = gerror.Wrap(err, "更新Foo状态失败，请稍后重试！")
	}
	return
}

// ================================================================
// 示例: Hgorm ORM 辅助工具 (internal/library/hgorm/)
// ================================================================

// ---------- handler/handler.go: Model 预处理选项 ----------
package handler

import (
	"github.com/gogf/gf/v2/database/gdb"
)

type Option struct {
	FilterAuth   bool // 过滤权限
	ForceCache   bool // 强制缓存
	FilterTenant bool // 过滤多租户权限
}

var DefaultOption = &Option{FilterAuth: true}

func Model(m *gdb.Model, opt ...*Option) *gdb.Model {
	var option *Option
	if len(opt) > 0 {
		option = opt[0]
	} else {
		option = DefaultOption
	}
	if option.FilterAuth {
		m = m.Handler(FilterAuth)
	}
	if option.ForceCache {
		m = m.Handler(ForceCache)
	}
	if option.FilterTenant {
		m = m.Handler(FilterTenant)
	}
	return m
}


// ---------- dao.go: IsUnique / LeftJoin / GetPkField 等 ----------
package hgorm

// IsUnique 唯一性检查（新增/编辑时判断字段是否重复）
func IsUnique(ctx context.Context, dao daoInstance, where g.Map, message string, pkId ...interface{}) error {
	m := dao.Ctx(ctx).Where(where)
	if len(pkId) > 0 {
		field, _ := GetPkField(ctx, dao)
		m = m.WhereNot(field, pkId[0])
	}
	exist, _ := m.Exist()
	if exist {
		return gerror.New(message)
	}
	return nil
}

// LeftJoin 关联表左连接（带自动别名和字段匹配）
func LeftJoin(m *gdb.Model, masterTable, masterField, joinTable, alias, onField string) *gdb.Model {
	return m.LeftJoin(GenJoinOnRelation(masterTable, masterField, joinTable, alias, onField)...)
}

// GetPkField 获取主键字段名
func GetPkField(ctx context.Context, dao daoInstance) (string, error) {
	fields, err := dao.Ctx(ctx).TableFields(dao.Table())
	for _, field := range fields {
		if strings.ToUpper(field.Key) == "PRI" {
			return field.Name, nil
		}
	}
	return "", gerror.New("no primary key")
}

// FilterKeywordsWithOr 多条件关键词OR查询
func FilterKeywordsWithOr(m *gdb.Model, filterColumns map[string]string, keyword string) *gdb.Model {
	var conditions []string
	var args []interface{}
	for col, operator := range filterColumns {
		switch operator {
		case "LIKE":
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+keyword+"%")
		default:
			conditions = append(conditions, fmt.Sprintf("%s = ?", col))
			args = append(args, keyword)
		}
	}
	return m.Where(fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")), args...)
}

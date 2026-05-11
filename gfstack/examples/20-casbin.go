// ================================================================
// 示例: Casbin RBAC 权限 (internal/library/casbin/enforcer.go)
// 基于 Casbin 的 RBAC 权限控制：从数据库加载角色-菜单-权限策略
// ================================================================

package casbin

import (
	"context"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
)

var Enforcer *casbin.Enforcer

// InitEnforcer 初始化 Casbin 执行器
func InitEnforcer(ctx context.Context) {
	// 1. 创建数据库适配器
	link := getDbLink(ctx)
	a, err := NewAdapter(link.String())
	if err != nil {
		g.Log().Panicf(ctx, "casbin.NewAdapter err: %v", err)
		return
	}

	// 2. 加载策略模型（优先本地 > 打包资源）
	path := "manifest/config/casbin.conf"
	modelContent := gfile.GetContents(path)
	if len(modelContent) == 0 && !gres.IsEmpty() {
		modelContent = string(gres.GetContent(path))
	}
	if len(modelContent) == 0 {
		g.Log().Panicf(ctx, "casbin model file does not exist: %v", path)
	}

	m, err := model.NewModelFromString(modelContent)
	if err != nil {
		g.Log().Panicf(ctx, "casbin NewModelFromString err: %v", err)
	}

	// 3. 创建执行器
	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		g.Log().Panicf(ctx, "casbin NewEnforcer err: %v", err)
	}

	// 4. 加载权限策略
	loadPermissions(ctx)
}

// getDbLink 获取数据库连接配置
func getDbLink(ctx context.Context) *gvar.Var {
	link := g.Cfg().MustGet(ctx, "database.default")
	if !link.IsSlice() {
		return g.Cfg().MustGet(ctx, "database.default.link")
	}
	// 读写分离：取主库
	for _, v := range link.Array() {
		val := v.(map[string]interface{})
		if val["role"] == "master" {
			return gvar.New(val["link"])
		}
	}
	return gvar.New("database.default.0.link")
}

// loadPermissions 从数据库加载角色-菜单权限
func loadPermissions(ctx context.Context) {
	type Policy struct {
		Key         string `json:"key"`
		Permissions string `json:"permissions"`
	}
	var (
		rules   [][]string
		polices []*Policy
	)

	// 联表查询: admin_role → admin_role_menu → admin_menu
	q := func(alias, column string) string {
		return fmt.Sprintf("%s.%s", alias, column)
	}
	err := g.Model(gstr.Join([]string{dao.AdminRole.Table(), "r"}, " ")).
		LeftJoin(gstr.Join([]string{dao.AdminRoleMenu.Table(), "rm"}, " "), "r.id=rm.role_id").
		LeftJoin(gstr.Join([]string{dao.AdminMenu.Table(), "m"}, " "), "rm.menu_id=m.id").
		Fields(q("r", dao.AdminRole.Columns().Key), q("m", dao.AdminMenu.Columns().Permissions)).
		Where(q("r", dao.AdminRole.Columns().Status), consts.StatusEnabled).
		Where(q("m", dao.AdminMenu.Columns().Status), consts.StatusEnabled).
		WhereNot(q("m", dao.AdminMenu.Columns().Permissions), "").
		WhereNot(q("r", dao.AdminRole.Columns().Key), consts.SuperRoleKey).
		Scan(&polices)
	if err != nil {
		g.Log().Fatalf(ctx, "loadPermissions Scan err:%v", err)
		return
	}

	for _, policy := range polices {
		if strings.Contains(policy.Permissions, ",") {
			for _, perm := range strings.Split(policy.Permissions, ",") {
				rules = append(rules, []string{policy.Key, perm, "GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD"})
			}
		} else {
			rules = append(rules, []string{policy.Key, policy.Permissions, "GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD"})
		}
	}

	if _, err = Enforcer.AddPolicies(rules); err != nil {
		g.Log().Fatalf(ctx, "loadPermissions AddPolicies err:%v", err)
	}
}

// Refresh 刷新权限策略
func Refresh(ctx context.Context) error {
	policy, _ := Enforcer.GetPolicy()
	if len(policy) > 0 {
		Enforcer.RemovePolicies(policy)
	}
	loadPermissions(ctx)
	return nil
}

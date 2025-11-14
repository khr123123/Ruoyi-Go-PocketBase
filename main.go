// main.go
package main

import (
	"Ruoyi-Go-PocketBase/models"
	"Ruoyi-Go-PocketBase/service"
	"fmt"
	"github.com/pocketbase/dbx"
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 获取用户权限菜单树
		se.Router.GET("/api/system/menu/tree", func(e *core.RequestEvent) error {
			info, err := e.RequestInfo()
			// 获取用户角色
			roleIdsRaw := info.Auth.Get("role")
			roleIds := parseRoleIds(roleIdsRaw)
			// 调用服务获取菜单树
			menuService := services.NewMenuService(app)
			menuTree, err := menuService.GetUserMenuTree(roleIds)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]any{
					"code":    500,
					"message": "查询菜单失败",
					"error":   err.Error(),
				})
			}
			return e.JSON(http.StatusOK, map[string]any{
				"code":    200,
				"message": "success",
				"data":    menuTree,
			})
		})

		// 获取所有菜单树（管理员使用）
		se.Router.GET("/api/system/menu/tree/all", func(e *core.RequestEvent) error {
			info, err := e.RequestInfo()
			// 可以在这里添加管理员权限检查
			_ = info // 使用 info 进行权限验证
			menuService := services.NewMenuService(app)
			menuTree, err := menuService.GetAllMenuTree()
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]any{
					"code":    500,
					"message": "查询菜单失败",
					"error":   err.Error(),
				})
			}
			return e.JSON(http.StatusOK, map[string]any{
				"code":    200,
				"message": "success",
				"data":    menuTree,
			})
		})

		// 原有的 settings 接口
		se.Router.GET("/api/myapp/settings", func(e *core.RequestEvent) error {
			info, _ := e.RequestInfo()
			roleIdsRaw := info.Auth.Get("role")
			roleIds := parseRoleIds(roleIdsRaw)

			roles := []models.SysRole{}
			err := app.DB().
				Select("*").
				From("sys_role").
				Where(dbx.In("id", roleIds)).
				All(&roles)

			fmt.Println(roles, err)

			return e.JSON(http.StatusOK, map[string]any{
				"roleIds": roleIds,
				"roles":   roles,
			})
		})

		se.Router.GET("/", apis.Static(os.DirFS("./pb_public/index.html"), false))
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// parseRoleIds 解析角色ID
func parseRoleIds(roleIdsRaw any) []string {
	roleIds := []string{}

	switch v := roleIdsRaw.(type) {
	case string:
		if v != "" {
			roleIds = []string{v}
		}
	case []string:
		roleIds = v
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok {
				roleIds = append(roleIds, s)
			}
		}
	default:
		fmt.Println("unknown role type:", v)
	}

	return roleIds
}

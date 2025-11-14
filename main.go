// main.go
package main

import (
	"Ruoyi-Go-PocketBase/service"
	"fmt"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"net/http"
)

func main() {
	app := pocketbase.New()
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 获取用户权限菜单树
		se.Router.GET("/api/getUserRouter", func(e *core.RequestEvent) error {
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
		se.Router.GET("/api/getAllMenuTree", func(e *core.RequestEvent) error {
			// 调用服务获取菜单树
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
		//se.Router.GET("/", apis.Static(os.DirFS("./pb_public/index.html"), false))
		return se.Next()
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// parseRoleIds 解析角色ID
func parseRoleIds(roleIdsRaw any) []string {
	var roleIds []string
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

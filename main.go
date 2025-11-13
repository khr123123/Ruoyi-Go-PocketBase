package main

import (
	"Ruoyi-Go-PocketBase/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/myapp/settings", func(e *core.RequestEvent) error {
			user := e.Auth
			ba := e.Auth.ExpandedAll("role")
			fmt.Println(ba)
			var roleIds []string
			switch v := user.Get("role").(type) {
			case []any:
				for _, r := range v {
					if s, ok := r.(string); ok {
						roleIds = append(roleIds, s)
					}
				}
			case string:
				roleIds = append(roleIds, v)
			}
			fmt.Println(roleIds)
			if len(roleIds) == 0 {
				return e.JSON(http.StatusOK, map[string]any{
					"userId":  user.Id,
					"roleIds": []string{},
					"menus":   []any{},
				})
			}
			var roles []models.SysRole
			err := e.App.DB().
				Select("id", "permission").
				From("sys_role").
				AndWhere(dbx.In("id", roleIds)).
				All(&roles)
			if err != nil {
				return err
			}
			var menuIds []string
			for _, role := range roles {
				menuIds = append(menuIds, role.Permission...)
			}
			menuSet := map[string]bool{}
			for _, id := range menuIds {
				menuSet[id] = true
			}
			var menus []models.SysMenu
			err = e.App.DB().
				Select("id", "menuName", "url").
				From("sys_menu").
				AndWhere(dbx.In("id", menuSet)).
				All(&menus)
			if err != nil {
				return err
			}
			return e.JSON(http.StatusOK, map[string]any{
				"userId":  user.Id,
				"roleIds": roleIds,
				"menus":   menus,
			})
		})
		se.Router.GET("/", apis.Static(os.DirFS("./pb_public/index.html"), false))
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

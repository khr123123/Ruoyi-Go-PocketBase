// main.go
package main

import (
	"Ruoyi-Go-PocketBase/models"
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
	app.OnRecordViewRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		// e.App
		// e.Collection
		// e.Record
		// and all RequestEvent fields...
		return e.Next()
	})
	app.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		// e.App
		// e.Collection
		// e.Record
		// and all RequestEvent fields...
		return e.Next()
	})
	// fires for every collection
	app.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		// e.App
		// e.Collection
		// e.Record
		// and all RequestEvent fields...
		return e.Next()
	})
	app.OnRecordDeleteRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		// e.App
		// e.Collection
		// e.Record
		// and all RequestEvent fields...
		return e.Next()
	})
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 获取用户权限菜单树
		se.Router.GET("/api/getUserRouter", func(e *core.RequestEvent) error {
			info, err := e.RequestInfo()
			// 获取用户角色
			roleIdsRaw := info.Auth.Get("role")
			// 情况1：roleIdsRaw 本身就是 []string
			roleIds, _ := roleIdsRaw.([]string)
			// 调用服务获取菜单树
			menuService := NewMenuService(app)
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
		}).Bind(apis.RequireAuth())
		se.Router.GET("/api/getAllMenuTree", func(e *core.RequestEvent) error {
			// 调用服务获取菜单树
			menuService := NewMenuService(app)
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
		se.Router.GET("/", apis.Static(os.DirFS("./pb_public/index.html"), false))
		return se.Next()
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// HasPermission 检查用户是否有指定权限
func HasPermission(user *models.User, required string, db *dbx.DB) bool {
	if user == nil {
		return false
	}
	// 根据用户角色查出所有权限
	var roles []models.SysRole
	roleIds := make([]interface{}, len(user.Role))
	for i, r := range user.Role {
		roleIds[i] = r
	}
	db.Select("*").From("sys_role").Where(dbx.In("id", roleIds...)).All(&roles)
	// 构建权限集合
	permissionSet := map[string]struct{}{}
	for _, role := range roles {
		for _, p := range role.Permission {
			permissionSet[p] = struct{}{}
		}
	}
	// 判断
	_, ok := permissionSet[required]
	return ok
}

type MenuService struct {
	app *pocketbase.PocketBase
}

func NewMenuService(app *pocketbase.PocketBase) *MenuService {
	return &MenuService{app: app}
}

// GetUserMenuTree 获取用户权限菜单树
func (s *MenuService) GetUserMenuTree(roleIds []string) ([]*models.MenuVO, error) {
	if len(roleIds) == 0 {
		return []*models.MenuVO{}, nil
	}
	permissions, err := s.getPermissionsByRoleIds(roleIds)
	if err != nil {
		return nil, err
	}
	menus, err := s.getMenusByPermissions(permissions)
	if err != nil {
		return nil, err
	}
	menuVOs := makeMenuVOs(menus)
	tree := s.buildMenuTree(menuVOs)
	return tree, nil
}

// GetAllMenuTree 获取所有菜单树（管理员用）
func (s *MenuService) GetAllMenuTree() ([]*models.MenuVO, error) {
	menus := []models.SysMenu{}
	if err := s.app.DB().Select("*").From("sys_menu").OrderBy("orderNum ASC", "created ASC").All(&menus); err != nil {
		return nil, err
	}
	menuVOs := makeMenuVOs(menus)
	return s.buildMenuTree(menuVOs), nil
}

// getPermissionsByRoleIds 根据角色ID获取权限集合
func (s *MenuService) getPermissionsByRoleIds(roleIds []string) ([]string, error) {
	roleIdsInterface := stringSliceToInterface(roleIds)
	var roles []models.SysRole
	if err := s.app.DB().Select("*").From("sys_role").Where(dbx.In("id", roleIdsInterface...)).All(&roles); err != nil {
		return nil, err
	}
	permissionSet := map[string]struct{}{}
	for _, r := range roles {
		for _, p := range r.Permission {
			permissionSet[p] = struct{}{}
		}
	}
	permissions := make([]string, 0, len(permissionSet))
	for p := range permissionSet {
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// getMenusByPermissions 根据权限列表查询菜单
func (s *MenuService) getMenusByPermissions(permissions []string) ([]models.SysMenu, error) {
	if len(permissions) == 0 {
		return []models.SysMenu{}, nil
	}
	permissionInterface := stringSliceToInterface(permissions)
	var menus []models.SysMenu
	if err := s.app.DB().
		Select("*").
		From("sys_menu").
		OrderBy("orderNum ASC", "created ASC").
		Where(dbx.In("id", permissionInterface...)).
		All(&menus); err != nil {
		return nil, err
	}
	return menus, nil
}

// buildMenuTree 构建树形菜单
func (s *MenuService) buildMenuTree(menus []*models.MenuVO) []*models.MenuVO {
	menuMap := make(map[string]*models.MenuVO, len(menus))
	var roots []*models.MenuVO
	for _, menu := range menus {
		menu.Children = []*models.MenuVO{}
		menuMap[menu.ID] = menu
	}
	for _, menu := range menus {
		if menu.ParentId == "" || menu.ParentId == "0" || menu.ParentId == "null" {
			roots = append(roots, menu)
			continue
		}
		if parent, ok := menuMap[menu.ParentId]; ok {
			parent.Children = append(parent.Children, menu)
		}
	}
	return roots
}

// makeMenuVOs 将 SysMenu 转换为 MenuVO
func makeMenuVOs(menus []models.SysMenu) []*models.MenuVO {
	menuVOs := make([]*models.MenuVO, len(menus))
	for i, m := range menus {
		menuVOs[i] = &models.MenuVO{
			ID:         m.ID,
			Icon:       getStringValue(m.Icon),
			MenuName:   m.MenuName,
			MenuType:   m.MenuType,
			OrderNum:   m.OrderNum,
			ParentId:   getStringValue(m.ParentId),
			Permission: getStringValue(m.Permission),
			URL:        getStringValue(m.URL),
			Children:   []*models.MenuVO{},
		}
	}
	return menuVOs
}

// 辅助函数：[]string → []interface{}
func stringSliceToInterface(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// getStringValue 从指针获取字符串值
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

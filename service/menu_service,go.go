// services/menu_service.go
package services

import (
	"Ruoyi-Go-PocketBase/models"
	"fmt"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

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
	// 1. 根据角色ID查询权限列表
	var roles []models.SysRole
	// 将 []string 转换为 []interface{}
	roleIdsInterface := make([]interface{}, len(roleIds))
	for i, v := range roleIds {
		roleIdsInterface[i] = v
	}
	err := s.app.DB().
		Select("*").
		From("sys_role").
		Where(dbx.In("id", roleIdsInterface...)).
		All(&roles)
	if err != nil {
		return nil, err
	}
	permissionSet := map[string]struct{}{}
	for _, r := range roles {
		for _, p := range r.Permission {
			permissionSet[p] = struct{}{}
		}
	}
	permissionInterface := make([]interface{}, 0, len(permissionSet))
	for p := range permissionSet {
		permissionInterface = append(permissionInterface, p)
	}
	fmt.Print(permissionInterface)
	// 3. 查询所有菜单
	var menus []models.SysMenu
	err = s.app.DB().
		Select("*").
		From("sys_menu").
		OrderBy("orderNum ASC", "created ASC").
		Where(dbx.In("id", permissionInterface...)).
		All(&menus)
	fmt.Println(menus)
	if err != nil {
		return nil, err
	}
	// 4. 过滤有权限的菜单并转换为 VO
	menuVOMap := make(map[string]*models.MenuVO)
	var filteredMenus []*models.MenuVO

	for _, menu := range menus {
		// 检查权限（如果菜单没有设置权限标识，则默认可见）
		permission := getStringValue(menu.Permission)
		menuVO := &models.MenuVO{
			ID:         menu.ID,
			Icon:       getStringValue(menu.Icon),
			MenuName:   menu.MenuName,
			MenuType:   menu.MenuType,
			OrderNum:   menu.OrderNum,
			ParentId:   getStringValue(menu.ParentId),
			Permission: permission,
			URL:        getStringValue(menu.URL),
			Children:   []*models.MenuVO{},
		}

		menuVOMap[menu.ID] = menuVO
		filteredMenus = append(filteredMenus, menuVO)
	}
	// 5. 构建树形结构
	tree := s.buildMenuTree(filteredMenus, menuVOMap)
	return tree, nil
}

// buildMenuTree 构建树形菜单
func (s *MenuService) buildMenuTree(menus []*models.MenuVO, menuMap map[string]*models.MenuVO) []*models.MenuVO {
	var roots []*models.MenuVO
	for _, menu := range menus {
		fmt.Printf("MenuID=%s ParentID=%s\n", menu.ID, menu.ParentId)
	}

	for _, menu := range menus {
		p := menu.ParentId
		// 根节点：ParentId 为 "" 或 "0" 或 "null"
		if p == "" || p == "0" || p == "null" {
			roots = append(roots, menu)
			continue
		}
		// 子节点
		if parent, exists := menuMap[p]; exists {
			parent.Children = append(parent.Children, menu)
		}
	}
	return roots
}

// GetAllMenuTree 获取所有菜单树（不过滤权限，用于管理员）
func (s *MenuService) GetAllMenuTree() ([]*models.MenuVO, error) {
	menus := []models.SysMenu{}
	err := s.app.DB().
		Select("*").
		From("sys_menu").
		OrderBy("orderNum ASC", "created ASC").
		All(&menus)

	if err != nil {
		return nil, err
	}

	menuVOMap := make(map[string]*models.MenuVO)
	var allMenus []*models.MenuVO

	for _, menu := range menus {
		menuVO := &models.MenuVO{
			ID:         menu.ID,
			Icon:       getStringValue(menu.Icon),
			MenuName:   menu.MenuName,
			MenuType:   menu.MenuType,
			OrderNum:   menu.OrderNum,
			ParentId:   getStringValue(menu.ParentId),
			Permission: getStringValue(menu.Permission),
			URL:        getStringValue(menu.URL),
			Children:   []*models.MenuVO{},
		}

		menuVOMap[menu.ID] = menuVO
		allMenus = append(allMenus, menuVO)
	}

	tree := s.buildMenuTree(allMenus, menuVOMap)

	return tree, nil
}

// getStringValue 辅助函数：从指针获取字符串值，如果为 nil 返回空字符串
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

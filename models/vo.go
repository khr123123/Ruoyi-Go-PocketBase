package models

// models/vo.go

type MenuVO struct {
	ID         string    `json:"id"`
	Icon       string    `json:"icon"`
	MenuName   string    `json:"menuName"`
	MenuType   string    `json:"menuType"`
	OrderNum   *int      `json:"orderNum"`
	ParentId   string    `json:"parentId"`
	Permission string    `json:"permission"`
	URL        string    `json:"url"`
	Children   []*MenuVO `json:"children,omitempty"`
}

// models/models.go
package models

import "github.com/pocketbase/pocketbase/tools/types"

type User struct {
	ID              string                  `json:"id" db:"id" gorm:"primaryKey;type:TEXT;default:'r'||lower(hex(randomblob(7)))"`
	Avatar          string                  `json:"avatar" db:"avatar" gorm:"type:TEXT;not null;default:''"`
	Created         string                  `json:"created" db:"created" gorm:"type:TEXT;not null;default:''"`
	Updated         string                  `json:"updated" db:"updated" gorm:"type:TEXT;not null;default:''"`
	Email           string                  `json:"email" db:"email" gorm:"type:TEXT;not null;default:'';uniqueIndex:idx_email__pb_users_auth_"`
	EmailVisibility bool                    `json:"emailVisibility" db:"emailVisibility" gorm:"type:BOOLEAN;not null;default:false"`
	Name            string                  `json:"name" db:"name" gorm:"type:TEXT;not null;default:''"`
	Password        string                  `json:"password" db:"password" gorm:"type:TEXT;not null;default:''"`
	TokenKey        string                  `json:"tokenKey" db:"tokenKey" gorm:"type:TEXT;not null;default:'';uniqueIndex:idx_tokenKey__pb_users_auth_"`
	Verified        bool                    `json:"verified" db:"verified" gorm:"type:BOOLEAN;not null;default:false"`
	Role            types.JSONArray[string] `json:"role" db:"role" gorm:"type:JSON;not null;default:'[]'"` // PocketBase JSON 数组
}

type SysRole struct {
	ID         string                  `json:"id" db:"id" gorm:"primaryKey;type:TEXT;default:'r'||lower(hex(randomblob(7)))"`
	Created    string                  `json:"created" db:"created" gorm:"type:TEXT;not null;default:''"`
	Updated    string                  `json:"updated" db:"updated" gorm:"type:TEXT;not null;default:''"`
	RoleName   string                  `json:"roleName" db:"roleName" gorm:"type:TEXT;not null;default:''"`
	Permission types.JSONArray[string] `json:"permission" db:"permission" gorm:"type:JSON;not null;default:'[]'"` // PocketBase JSON 数组
}

type SysMenu struct {
	ID         string  `json:"id" db:"id" gorm:"primaryKey;type:TEXT;default:'r'||lower(hex(randomblob(7)))"`
	Created    string  `json:"created" db:"created" gorm:"type:TEXT;not null;default:''"`
	Updated    string  `json:"updated" db:"updated" gorm:"type:TEXT;not null;default:''"`
	Icon       *string `json:"icon" db:"icon" gorm:"type:TEXT;default:null"`
	MenuName   string  `json:"menuName" db:"menuName" gorm:"type:TEXT;not null;default:''"`
	MenuType   string  `json:"menuType" db:"menuType" gorm:"type:TEXT;not null;default:''"`
	OrderNum   *int    `json:"orderNum" db:"orderNum" gorm:"type:NUMERIC;default:null"`
	ParentId   *string `json:"parentId" db:"parentId" gorm:"type:TEXT;default:null"`
	Permission *string `json:"permission" db:"permission" gorm:"type:TEXT;default:null"`
	URL        *string `json:"url" db:"url" gorm:"type:TEXT;default:null"`
}

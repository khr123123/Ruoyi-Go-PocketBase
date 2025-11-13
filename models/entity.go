package models

type User struct {
	ID              string   `json:"id" gorm:"primaryKey;type:TEXT;default:'r'||lower(hex(randomblob(7)))"`
	Avatar          string   `json:"avatar" gorm:"type:TEXT;not null;default:''"`
	Created         string   `json:"created" gorm:"type:TEXT;not null;default:''"`
	Updated         string   `json:"updated" gorm:"type:TEXT;not null;default:''"`
	Email           string   `json:"email" gorm:"type:TEXT;not null;default:'';uniqueIndex:idx_email__pb_users_auth_"`
	EmailVisibility bool     `json:"emailVisibility" gorm:"type:BOOLEAN;not null;default:false"`
	Name            string   `json:"name" gorm:"type:TEXT;not null;default:''"`
	Password        string   `json:"password" gorm:"type:TEXT;not null;default:''"`
	TokenKey        string   `json:"tokenKey" gorm:"type:TEXT;not null;default:'';uniqueIndex:idx_tokenKey__pb_users_auth_"`
	Verified        bool     `json:"verified" gorm:"type:BOOLEAN;not null;default:false"`
	Role            []string `json:"role" gorm:"type:JSON;not null;default:'[]'"`
}

type SysRole struct {
	ID         string   `json:"id" gorm:"primaryKey;type:TEXT;default:'r'||lower(hex(randomblob(7)))"`
	Created    string   `json:"created" gorm:"type:TEXT;not null;default:''"`
	Updated    string   `json:"updated" gorm:"type:TEXT;not null;default:''"`
	RoleName   string   `json:"roleName" gorm:"type:TEXT;not null;default:''"`
	Permission []string `json:"permission" gorm:"type:JSON;not null;default:'[]'"`
}

type SysMenu struct {
	ID         string `json:"id" gorm:"primaryKey;type:TEXT;default:'r'||lower(hex(randomblob(7)))"`
	Created    string `json:"created" gorm:"type:TEXT;not null;default:''"`
	Updated    string `json:"updated" gorm:"type:TEXT;not null;default:''"`
	Icon       string `json:"icon" gorm:"type:TEXT;default:''"`
	MenuName   string `json:"menuName" gorm:"type:TEXT;not null;default:''"`
	MenuType   string `json:"menuType" gorm:"type:TEXT;not null;default:''"`
	OrderNum   int    `json:"orderNum" gorm:"type:NUMERIC;default:0"`
	ParentId   string `json:"parentId" gorm:"type:TEXT;default:''"` // 注意原表是 parendId，可能是拼写错误
	Permission string `json:"permission" gorm:"type:TEXT;default:''"`
	URL        string `json:"url" gorm:"type:TEXT;default:''"`
}

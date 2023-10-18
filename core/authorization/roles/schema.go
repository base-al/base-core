package roles

import (
	"base-core/core/orgs"
	"time"
)

type Role struct {
	ID          int `gorm:"primaryKey"`
	OrgID       *int
	Name        string
	Description *string
	UserOrgRole []orgs.UserOrgRole `gorm:"foreignKey:RoleID"`
	Permissions []Permission       `gorm:"many2many:role_permissions"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type Permission struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	HTTPMethod  string
	Path        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

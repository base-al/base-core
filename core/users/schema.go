package users

import (
	"base-core/core/orgs"
	"time"
)

type User struct {
	ID            int    `gorm:"primaryKey"`
	Email         string `gorm:"unique"`
	Username      string `gorm:"unique"`
	Password      string
	FirstName     string
	LastName      string
	Active        bool
	UserOrgRole   []orgs.UserOrgRole `gorm:"foreignKey:UserID"`
	VerifiedEmail bool
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}

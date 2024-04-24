package google

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const OAuthAccountTableName = "oauth_accounts"

type OAuthAccount struct {
	ID           int `gorm:"primaryKey"`
	UserID       int `gorm:"foreignKey:ID"`
	User         User
	Email        string
	Provider     string
	OAuthVersion string `gorm:"column:oauth_version"`
	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
	Active       bool
	Scopes       pq.StringArray  `gorm:"type:text[]"`
	ExtraData    json.RawMessage `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

func (OAuthAccount) TableName() string {
	return OAuthAccountTableName
}

// Note: Do not change here
// Reference: base/core/users/schema.go
type User struct {
	gorm.Model
	Email         string  `gorm:"unique"`
	Username      *string `gorm:"unique"`
	Password      string
	FirstName     string
	LastName      string
	Active        bool
	Phone         string
	VerifiedEmail bool
}

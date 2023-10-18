package orgs

import (
	"base-core/core/subscriptions"
	"time"
)

type Org struct {
	ID              int           `gorm:"primaryKey"`
	Name            string        `gorm:"not null"`
	TradeName       string        `gorm:"not null"`
	Size            string        `gorm:"not null"`
	Slug            string        `gorm:"unique;not null"`
	UserCompanyRole []UserOrgRole `gorm:"foreignKey:OrgID"`
	UIN             string
	VAT             string
	SubscriptionID  int
	Subscription    subscriptions.Subscription
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

type OrgSettings struct {
	ID          int `gorm:"primaryKey"`
	OrgID       int
	Org         Org
	WebsiteURL  string
	Country     string
	Address     string
	City        string
	ZIP         string
	Domain      string
	Email       string
	PhoneNumber string
	Industry    string
	Theme       string
	TermsURL    string
	LogoImgURL  string
	LinkedinURL string
	FacebookURL string
	XURL        string
	GithubURL   string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type UserOrgRole struct {
	UserID int
	OrgID  int
	RoleID int
	Status string
}

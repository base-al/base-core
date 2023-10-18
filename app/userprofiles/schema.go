package userprofiles

import (
	"base-core/core/users"
	"time"
)

type Profile struct {
	ID           int `gorm:"primaryKey"`
	UserID       int
	User         users.User
	Headline     string
	AvatarImgUrl string
	Bio          string
	Country      string
	Address      string
	LinkedinURL  string
	Website      string
	XURL         string
	Department   string
	Skills       string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

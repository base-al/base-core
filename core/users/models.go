package users

import (
	"time"
)

type SignupRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type SignupResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type SignupVerifyRequest struct {
	Token string `json:"-"`
}

type StatusResponse struct {
	Status bool `json:"status"`
}

type UserRequest struct {
	UserID    int    `json:"-"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type IDResponse struct {
	ID int `json:"id"`
}

type DonwloadAvatarImgRequest struct {
	UserID int `json:"-"`
}

type DownloadAvatarImgResponse struct {
	AvatarImgURL string `json:"avatarImgUrl"`
}

type UserResponse struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Status       string    `json:"status"`
	AvatarImgKey string    `json:"avatarImgUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PasswordUpdateRequest struct {
	UserID             int    `json:"-"`
	Mode               string `json:"mode"`
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type EmailUpdateRequest struct {
	UserID   int    `json:"-"`
	NewEmail string `json:"newEmail"`
}

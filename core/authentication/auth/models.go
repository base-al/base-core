package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AvatarImgUrl string `json:"avatarImgUrl"`
}

type OrgNameRole struct {
	OrgID  int    `json:"orgId"`
	RoleID int    `json:"roleId"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

type LoginResponse struct {
	UserData *UserData     `json:"userData"`
	Orgs     []OrgNameRole `json:"orgs"`
	Token    string        `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token              string `json:"-"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type StatusReply struct {
	Status bool `json:"status"`
}

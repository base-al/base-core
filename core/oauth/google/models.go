package oauth

import "time"

type OAuthResponse struct {
	RedirectURL string `json:"redirectUrl"`
}

type OAuthCallbackRequest struct {
	Code string `json:"code"`
}

type CallbackResponse struct {
	Provider string
	Token    string
}

type FindUserOAuthAccount struct {
	UserID int `json:"-"`
}

type UserOAuthAccount struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserAccountResponse struct {
	OAuthAccounts []*UserOAuthAccount `json:"oAuthAccounts"`
	PasswordSet   bool                `json:"passwordSet"`
}

type GoogleAccountConnectRequest struct {
	BaseUserID int `json:"-"`
}

type GoogleAccountConnectCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type StatusResponse struct {
	Status bool `json:"status"`
}

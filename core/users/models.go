package users

type SignupUserRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type SignupUserResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type SignupUserVerifyRequest struct {
	Token string `json:"-"`
}

type StatusResponse struct {
	Status bool `json:"status"`
}

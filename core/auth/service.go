package auth

import (
	"base/helper"
	"base/s3"
	"context"
	"fmt"
	"strings"

	"time"

	// companysvc "base/app/companies"

	oauth "base/core/oauth/google"
	userService "base/core/users"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mattevans/postmark-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type authAPI struct {
	db             *gorm.DB
	secretKey      string
	postmarkClient *postmark.Client
	s3c            *s3.S3Client
	logger         log.AllLogger
	uiAppUrl       string
	oauth2Config   *oauth2.Config
	sendFromEmail  string
}

type AuthAPI interface {
	Signup(req *SignupRequest) (res *SignupResponse, err error)
	SignupVerify(req *SignupVerifyRequest) (res *StatusResponse, err error)
	Login(req *LoginRequest) (res *LoginResponse, err error)
	ForgotPassword(req *ForgotPasswordRequest) (res *StatusResponse, err error)
	ResetPassword(req *ResetPasswordRequest) (res *StatusResponse, err error)
	OAuthLogin(req *OAuthLoginRequest) (res *LoginResponse, err error)
}

func NewAuthAPI(db *gorm.DB, sk string, pmc *postmark.Client, s3c *s3.S3Client, uiAppUrl string, logger log.AllLogger, oauth2Config *oauth2.Config, sendFromEmail string) AuthAPI {
	return &authAPI{
		db:             db,
		secretKey:      sk,
		postmarkClient: pmc,
		s3c:            s3c,
		logger:         logger,
		uiAppUrl:       uiAppUrl,
		oauth2Config:   oauth2Config,
		sendFromEmail:  sendFromEmail,
	}
}

// @Summary      	Signup
// @Description	Validates email, username, first name, last name, password checks if email exists, if not creates new user and sends email with verification link.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			SignupRequest	body		SignupRequest	true	"SignupRequest"
// @Success			200					{object}	SignupResponse
// @Router			/auth/signup	[POST]
func (s authAPI) Signup(req *SignupRequest) (res *SignupResponse, err error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Password = strings.TrimSpace(req.Password)
	req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)

	if req.Email == "" {
		return nil, fmt.Errorf("Email is required")
	}
	if !helper.ValidEmail(req.Email) {
		return nil, fmt.Errorf("Email not valid")
	}
	if req.Username == "" {
		return nil, fmt.Errorf("Username is required")
	}
	if req.FirstName == "" {
		return nil, fmt.Errorf("FirstName is required")
	}
	if req.LastName == "" {
		return nil, fmt.Errorf("LastName is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("Password is required")
	}
	if req.ConfirmPassword == "" {
		return nil, fmt.Errorf("ConfirmPassword is required")
	}
	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("Password are ConfirmPassword must match")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["email"] = req.Email
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	var user userService.User
	_ = s.db.Where("email = ?", req.Email).First(&user)
	if user.ID > 0 {
		if !user.VerifiedEmail {
			verifyLink := s.uiAppUrl + "/verify-signup/" + t
			// Send email to req email
			emailMsg := &postmark.Email{
				From:    s.sendFromEmail,
				To:      req.Email,
				Subject: "Base - Verify your registration",
				HTMLBody: fmt.Sprintf(`Hello from Base!<br/><br/>
					You have successfully signed up as new user in Base platform!<br/>
					Please continue by verifing your registration by linking the link below.<br/><br/>

					<a href='%s'>Verify you email</a><br/><br/>

					Thank you, <br/>
					Base Team
				`, verifyLink),
			}

			_, _, err = s.postmarkClient.Email.Send(emailMsg)
			if err != nil {
				s.logger.Errorf("func: Signup, operation: s.postmarkClient.Email.Send, err: %s", err.Error())
				return nil, err
			}

			return &SignupResponse{
				ID:     user.ID,
				Status: "pending",
			}, nil
		}
		return nil, fmt.Errorf("Email is already used")
	}

	_ = s.db.Where("username = ?", req.Username).First(&user)
	if user.ID > 0 {
		return nil, fmt.Errorf("username is already used")

	}

	user.Email = req.Email
	user.Username = &req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.VerifiedEmail = false
	user.Active = false

	pwh, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("func: Signup, operation: bcrypt.GenerateFromPassword, err: %s", err.Error())
		return nil, err
	}
	pwhs := string(pwh)
	user.Password = pwhs

	result := s.db.Omit("UpdatedAt").Create(&user)
	if result.Error != nil {
		s.logger.Errorf("func: Signup, operation: s.db.Omit('UpdatedAt').Create(&user), err: %s", result.Error.Error())
		return nil, result.Error
	}
	s.db.Model(userService.User{Email: req.Email}).First(&user)

	// Create directory in users.env bucket for this user
	usrsBckName := s.s3c.UsersBucketName()

	err = s.s3c.NewBucketFolder(usrsBckName, fmt.Sprint(user.ID))
	if err != nil {
		s.logger.Errorf("func: Signup, operation: s.s3c.NewBucketFolder, err: %s", err.Error())
		return nil, err
	}

	verifyLink := s.uiAppUrl + "/verify-signup/" + t
	// Send email to req email
	fmt.Println("s.sendFromEmail", s.sendFromEmail)
	emailMsg := &postmark.Email{
		From:    s.sendFromEmail,
		To:      req.Email,
		Subject: "Base -  Verify your registration",
		HTMLBody: fmt.Sprintf(`Hello from Base!<br/><br/>
			You have successfully signed up as new user in Base platform!<br/>
			Please continue by verifing your registration by linking the link below.<br/><br/>

			<a href='%s'>Verify you email</a><br/><br/>

			Thank you, <br/>
			Base Team
		`, verifyLink),
	}

	_, _, err = s.postmarkClient.Email.Send(emailMsg)
	if err != nil {
		s.logger.Errorf("func: Signup, operation: s.postmarkClient.Email.Send, err: %s", err.Error())
		return nil, err
	}

	return &SignupResponse{
		ID:     user.ID,
		Status: "pending",
	}, nil
}

// @Summary      	SignupVerify
// @Description	Validates token in param, if token parses valid then user will be verified and be updated in DB.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			token				path		string			true	"Token"
// @Success			200					{object}	StatusResponse
// @Router			/auth/signup/verify/{token}	[PUT]
func (s authAPI) SignupVerify(req *SignupVerifyRequest) (res *StatusResponse, err error) {
	req.Token = strings.TrimSpace(req.Token)
	if req.Token == "" {
		return nil, fmt.Errorf("Token is empty")
	}

	// Parse and validate token expiration
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		s.logger.Errorf("func: SignupVerify, operation: jwt.ParseWithClaims, err: %s", err.Error())
		return nil, err
	}
	email := fmt.Sprintf("%v", claims["email"])

	// Find user by email in DB
	var user userService.User
	s.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, helper.ErrNotFound
	}

	user.VerifiedEmail = true
	user.Active = true
	result := s.db.Save(&user)
	if result.Error != nil {
		s.logger.Errorf("func: SignupVerify, operation: s.db.Save(&user), err: %s", result.Error.Error())
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

// @Summary      	Login
// @Description		Validates email and password in request, check if user exists in DB if not throw 404 otherwise compare the request password with hash, then check if user is active, then finds relationships of user with orgs and then generates a JWT token, and returns UserData, Orgs, and Token in response.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			LoginRequest	body		LoginRequest	true	"LoginRequest"
// @Success			200				{object}	LoginResponse
// @Router			/auth/login			[POST]
func (s authAPI) Login(req *LoginRequest) (res *LoginResponse, err error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	// Validate email
	if req.Email == "" {
		return nil, fmt.Errorf("email is empty")
	}
	if !helper.ValidEmail(req.Email) {
		return nil, fmt.Errorf("email not valid")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is empty")
	}

	// Check user exists
	var user userService.User
	s.db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		fmt.Println("err", "email not found")
		return nil, helper.ErrNotFound
	}

	if !user.VerifiedEmail {
		return nil, fmt.Errorf("email not verified")
	}
	// check validity
	if !user.Active {
		return nil, fmt.Errorf("user not active")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, helper.ErrNotFound
	}

	// Generate token with expiration
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = user.ID

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// var userCompany companysvc.UserCompany
	// _ = s.db.Where("user_id = ?", user.ID).First(&userCompany)
	// claims["company"] = userCompany

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	userData := UserData{
		ID: user.ID,

		Email:     user.Email,
		Username:  helper.PointerValue(user.Username),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return &LoginResponse{
		UserData: &userData,

		Token: t,
	}, nil
}

// @Summary      	ForgotPassword
// @Description		Validates email if exists in DB, then send an email with verification link to email user.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			ForgotPasswordRequest	body		ForgotPasswordRequest	true	"ForgotPasswordRequest"
// @Success			200					{object}	StatusResponse
// @Router			/auth/forgot-password	[POST]
func (s authAPI) ForgotPassword(req *ForgotPasswordRequest) (res *StatusResponse, err error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	// Find user by email in DB
	var user userService.User
	s.db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		return nil, helper.ErrNotFound
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["email"] = req.Email
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	verifyLink := s.uiAppUrl + "/reset-password/" + t
	// Send email to req email
	emailMsg := &postmark.Email{
		From:    "info@base.al",
		To:      req.Email,
		Subject: "Base -  Reset your password",
		HTMLBody: fmt.Sprintf(`Hello from Base!<br/><br/>
			You have requested to reset you password!<br/>
			Please continue resetting password by linking the link below.<br/><br/>

			<a href='%s'>Reset your password</a><br/><br/>

			Thank you, <br/>
			Base Team
		`, verifyLink),
	}

	_, _, err = s.postmarkClient.Email.Send(emailMsg)
	if err != nil {
		return nil, err
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

// @Summary      	ResetPassword
// @Description		Validates token, new password, and confirm new password, checks if user exists in DB then it updates the password in DB.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			token				path		string			true	"Token"
// @Param			ResetPasswordRequest	body		ResetPasswordRequest	true	"ResetPasswordRequest"
// @Success			200					{object}	StatusResponse
// @Router			/auth/reset-password/{token}	[POST]
func (s authAPI) ResetPassword(req *ResetPasswordRequest) (res *StatusResponse, err error) {
	req.Token = strings.TrimSpace(req.Token)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.ConfirmNewPassword = strings.TrimSpace(req.ConfirmNewPassword)
	if req.Token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	if req.NewPassword == "" {
		return nil, fmt.Errorf("new Password is required")
	}
	if req.ConfirmNewPassword == "" {
		return nil, fmt.Errorf("confirm New Password is required")
	}
	if req.NewPassword != req.ConfirmNewPassword {
		return nil, fmt.Errorf("new Password are Confirm New Password must match")
	}

	// Parse and validate token expiration
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		fmt.Println("ParseWithClaims", err)
		return nil, err
	}
	email := fmt.Sprintf("%v", claims["email"])

	// Find user by email in DB
	var user userService.User
	s.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, helper.ErrNotFound
	}

	pwh, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	pwhs := string(pwh)
	user.Password = pwhs

	result := s.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

func (s authAPI) OAuthLogin(req *OAuthLoginRequest) (res *LoginResponse, err error) {
	req.Provider = strings.TrimSpace(req.Provider)
	req.Token = strings.TrimSpace(req.Token)
	if req.Provider == "" {
		return nil, fmt.Errorf("provider is empty")
	}
	if req.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	switch req.Provider {
	case oauth.GoogleProviderName:
		// Parse and validate token expiration
		oauthclaims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(req.Token, oauthclaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKey), nil
		})
		if err != nil {
			s.logger.Errorf("func: Login, operation: jwt.ParseWithClaims, err: %s", err.Error())
			return nil, err
		}
		// reqEmail := fmt.Sprintf("%v", oauthclaims["email"])
		oAuthId := fmt.Sprintf("%v", oauthclaims["oid"])
		providerName := fmt.Sprintf("%v", oauthclaims["provider"])

		googleEmailScope := "https://www.googleapis.com/auth/userinfo.email"
		var oAuthAcc oauth.OAuthAccount
		result := s.db.Where("id = ? AND provider = ? AND ? = ANY(scopes)", oAuthId, providerName, googleEmailScope).First(&oAuthAcc)
		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				return nil, fmt.Errorf("account not connected, login with your email and password to connect with google.")
			}
			return nil, result.Error
		}
		// hasEmailScope
		// for scope := range oAuthAcc.Scopes {

		// }

		// TODO: Check user exists (handle when user updates)
		// var user userService.User
		// s.db.Where("email = ?", reqEmail).First(&user)
		// if user.ID == 0 {
		// 	fmt.Println("err", "email not found")
		// 	return nil, helper.ErrNotFound
		// }

		var user userService.User
		s.db.Where("id = ?", oAuthAcc.UserID).First(&user)
		if user.ID == 0 {
			return nil, helper.ErrNotFound
		}

		// check user profile for avatar img

		// check validity
		if !user.Active {
			return nil, fmt.Errorf("user not active")
		}

		if !user.VerifiedEmail {
			return nil, fmt.Errorf("email not verified")
		}

		// Update google access token if expired
		// oAuthAcc.AccessToken
		tok := &oauth2.Token{
			AccessToken:  oAuthAcc.AccessToken,
			TokenType:    "Bearer",
			RefreshToken: oAuthAcc.RefreshToken,
			Expiry:       oAuthAcc.TokenExpiry,
		}
		ctx := context.Background()
		if tok.Expiry.Before(time.Now()) {
			if oAuthAcc.RefreshToken == "" {
				return nil, fmt.Errorf("oauthacc refresh token missing")
			}
			newTok, err := s.oauth2Config.TokenSource(ctx, tok).Token()
			if err != nil {
				fmt.Printf("Failed to refresh access token: %v", err)
				return nil, err
			}
			tok = newTok
			// oAuthAcc.AccessToken = tok.AccessToken
			result := s.db.Model(&oAuthAcc).Update("access_token", tok.AccessToken)
			if result.Error != nil {
				return nil, fmt.Errorf("cannot update access token")
			}
		}

		// Generate token with expiration
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userid"] = user.ID

		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		// var userCompany companysvc.UserCompany
		// _ = s.db.Where("user_id = ?", user.ID).First(&userCompany)
		// claims["company"] = userCompany

		t, err := token.SignedString([]byte(s.secretKey))
		if err != nil {
			return nil, err
		}

		userData := UserData{
			ID:        user.ID,
			Email:     user.Email,
			Username:  helper.PointerValue(user.Username),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}

		return &LoginResponse{
			UserData: &userData,

			Token: t,
		}, nil

	default:
		return nil, fmt.Errorf("OAuth provider doesn't supported")
	}
}

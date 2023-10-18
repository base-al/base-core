package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"base-core/core/users"
	"base-core/helper"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"github.com/mattevans/postmark-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authAPI struct {
	db             *gorm.DB
	secretKey      string
	postmarkClient *postmark.Client
	logger         log.AllLogger
}

type AuthAPI interface {
	Login(req *LoginRequest) (res *LoginResponse, err error)
	ForgotPassword(req *ForgotPasswordRequest) (res *StatusReply, err error)
	ResetPassword(req *ResetPasswordRequest) (res *StatusReply, err error)
}

func NewAuthAPI(db *gorm.DB, sk string, pmc *postmark.Client, logger log.AllLogger) AuthAPI {
	return &authAPI{
		db:             db,
		secretKey:      sk,
		postmarkClient: pmc,
		logger:         logger,
	}
}

func (s authAPI) Login(req *LoginRequest) (res *LoginResponse, err error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	// Validate email
	if req.Email == "" {
		return nil, fmt.Errorf("Email is empty")
	}
	if !helper.ValidEmail(req.Email) {
		return nil, fmt.Errorf("Email not valid")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("Password is empty")
	}

	// Check user exists
	var user users.User
	s.db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		fmt.Println("err", "email not found")
		return nil, helper.ErrNotFound
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Println("err", "wrong pwd")
		return nil, helper.ErrNotFound
	}

	// check validity
	if !user.Active {
		return nil, fmt.Errorf("User not active")
	}

	if !user.VerifiedEmail {
		return nil, fmt.Errorf("Email not verified")
	}

	// Generate token with expiration
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = user.ID

	rows, err := s.db.Table("orgs").
		Select("id", "name", "slug", "role_id").
		Joins("LEFT JOIN user_org_roles ON user_org_roles.org_id=orgs.id").
		Where("user_org_roles.user_id = ?", user.ID).Order("id ASC").Rows()
	var orgNameRoles []OrgNameRole
	for rows.Next() {
		onr := OrgNameRole{}
		err := rows.Scan(&onr.OrgID, &onr.Name, &onr.Slug, &onr.RoleID)
		if err != nil {
			fmt.Println("err", err)
		}
		orgNameRoles = append(orgNameRoles, onr)
	}
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	userData := UserData{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return &LoginResponse{
		UserData: &userData,
		Orgs:     orgNameRoles,
		Token:    t,
	}, nil
}

func (s authAPI) ForgotPassword(req *ForgotPasswordRequest) (res *StatusReply, err error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, fmt.Errorf("Email is empty")
	}

	// Find user by email in DB
	var user users.User
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

	verifyLink := os.Getenv("APP_URL") + "/reset-password/" + t
	// Send email to req email
	emailMsg := &postmark.Email{
		From:    "info@communihub.co",
		To:      req.Email,
		Subject: "CommuniHub -  Reset your password",
		HTMLBody: fmt.Sprintf(`Hello from CommuniHub!<br/><br/>
			You have requested to reset you password!<br/>
			Please continue resetting password by linking the link below.<br/><br/>

			<a href='%s'>Reset your password</a><br/><br/>

			Thank you, <br/>
			CommuniHub Team
		`, verifyLink),
	}

	_, _, err = s.postmarkClient.Email.Send(emailMsg)
	if err != nil {
		return nil, err
	}

	return &StatusReply{
		Status: true,
	}, nil
}

func (s authAPI) ResetPassword(req *ResetPasswordRequest) (res *StatusReply, err error) {
	req.Token = strings.TrimSpace(req.Token)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.ConfirmNewPassword = strings.TrimSpace(req.ConfirmNewPassword)
	if req.Token == "" {
		return nil, fmt.Errorf("Token is empty")
	}
	if req.NewPassword == "" {
		return nil, fmt.Errorf("New Password is required")
	}
	if req.ConfirmNewPassword == "" {
		return nil, fmt.Errorf("Confirm New Password is required")
	}
	if req.NewPassword != req.ConfirmNewPassword {
		return nil, fmt.Errorf("New Password are Confirm New Password must match")
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
	var user users.User
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

	return &StatusReply{
		Status: true,
	}, nil
}

package users

import (
	"base-core/helper"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mattevans/postmark-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userAPI struct {
	db             *gorm.DB
	secretKey      string
	postmarkClient *postmark.Client
	logger         log.AllLogger
}

type UserAPI interface {
	Signup(req *SignupUserRequest) (res *SignupUserResponse, err error)
	SignupVerify(req *SignupUserVerifyRequest) (res *StatusResponse, err error)
}

func NewUserAPI(db *gorm.DB, sk string, pmc *postmark.Client, logger log.AllLogger) UserAPI {
	return &userAPI{
		db:             db,
		secretKey:      sk,
		postmarkClient: pmc,
		logger:         logger,
	}
}

// @Summary      	Signup
// @Description	Validates email, username, first name, last name, password checks if email exists, if not creates new user and sends email with verification link.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Param			SignupUserRequest	body		SignupUserRequest	true	"SignupUserRequest"
// @Success			200					{object}	SignupUserResponse
// @Router			/signup	[POST]
func (s userAPI) Signup(req *SignupUserRequest) (res *SignupUserResponse, err error) {
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

	var user User
	_ = s.db.Where("email = ?", req.Email).First(&user)
	if user.ID > 0 {
		if !user.VerifiedEmail {
			verifyLink := os.Getenv("APP_URL") + "/verify-signup/" + t
			// Send email to req email
			emailMsg := &postmark.Email{
				From:    "info@base.al",
				To:      req.Email,
				Subject: "Base -  Verify your registration",
				HTMLBody: fmt.Sprintf(`Hello from Base!<br/><br/>
					You have successfully signed up as new user in Base platform!<br/>
					Please continue by verifing your registration by linking the link below.<br/><br/>

					<a href='%s'>Verify you email</a><br/><br/>

					Thank you, <br/>
					CommuniHub Team
				`, verifyLink),
			}

			_, _, err = s.postmarkClient.Email.Send(emailMsg)
			if err != nil {
				return nil, err
			}

			return &SignupUserResponse{
				ID:     user.ID,
				Status: "pending",
			}, nil
		}
		return nil, fmt.Errorf("Email is already used")
	}

	_ = s.db.Where("username = ?", req.Username).First(&user)
	if user.ID > 0 {
		return nil, fmt.Errorf("Username is already used")
	}

	fmt.Println("token", t)

	// Save new profile and user record

	user.Email = req.Email
	user.Username = req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.VerifiedEmail = false
	user.Active = false

	pwh, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	pwhs := string(pwh)
	user.Password = pwhs

	result := s.db.Omit("UpdatedAt").Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	s.db.Model(User{Email: req.Email}).First(&user)

	// Create directory in users.env bucket for this user
	// usrsBckName := s.s3c.UsersBucketName()
	// err = s.s3c.NewBucketFolder(usrsBckName, fmt.Sprint(user.ID))
	// if err != nil {
	// 	return nil, err
	// }

	verifyLink := os.Getenv("APP_URL") + "/verify-signup/" + t
	// Send email to req email
	emailMsg := &postmark.Email{
		From:    "info@base.al",
		To:      req.Email,
		Subject: "Base -  Verify your registration",
		HTMLBody: fmt.Sprintf(`Hello from Base!<br/><br/>
			You have successfully signed up as new user in CommuniHub platform!<br/>
			Please continue by verifing your registration by linking the link below.<br/><br/>

			<a href='%s'>Verify you email</a><br/><br/>

			Thank you, <br/>
			CommuniHub Team
		`, verifyLink),
	}

	_, _, err = s.postmarkClient.Email.Send(emailMsg)
	if err != nil {
		return nil, err
	}

	return &SignupUserResponse{
		ID:     user.ID,
		Status: "pending",
	}, nil
}

// @Summary      	SignupVerify
// @Description	Validates token in param, if token parses valid then user will be verified and be updated in DB.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Param			token				path		string			true	"Token"
// @Success			200					{object}	StatusResponse
// @Router			/signup/verify/{token}	[GET]
func (s userAPI) SignupVerify(req *SignupUserVerifyRequest) (res *StatusResponse, err error) {
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
		fmt.Println("ParseWithClaims", err)
		return nil, err
	}
	email := fmt.Sprintf("%v", claims["email"])

	// Find user by email in DB
	var user User
	s.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, helper.ErrNotFound
	}

	user.VerifiedEmail = true
	user.Active = true
	result := s.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

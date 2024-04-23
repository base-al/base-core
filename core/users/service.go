package users

import (
	"fmt"
	"strings"

	helpers "github.com/base-al/base-core/helpers"
	"github.com/base-al/base-core/s3"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mattevans/postmark-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	UserStatusActive   = string("active")
	UserStatusInactive = string("inactive")
	UserStatusInvited  = string("invited")
	UserStatusPending  = string("pending")
	UserStatusReject   = string("rejected")

	modeChangePassword = string("change")
	modeSetPassword    = string("set")
)

type userAPI struct {
	db             *gorm.DB
	s3c            *s3.S3Client
	secretKey      string
	postmarkClient *postmark.Client
	uiAppUrl       string
	logger         log.AllLogger
	sendFromEmail  string
}

type UserAPI interface {
	PasswordUpdate(req *PasswordUpdateRequest) (res *StatusResponse, err error)
	EmailUpdate(req *EmailUpdateRequest) (res *StatusResponse, err error)
}

func NewUserAPI(db *gorm.DB, s3c *s3.S3Client, sk string, pmc *postmark.Client, uiAppUrl string, logger log.AllLogger, sendFromEmail string) UserAPI {
	return userAPI{
		db:             db,
		s3c:            s3c,
		secretKey:      sk,
		postmarkClient: pmc,
		uiAppUrl:       uiAppUrl,
		logger:         logger,
		sendFromEmail:  sendFromEmail,
	}
}

// @Summary      	PasswordUpdate
// @Description		Validates user id, mode, new password, confirm new password and(or) current password then will set or update password of the user.
// @Tags			Users
// @Produce			json
// @Param			Authorization						header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			PasswordUpdateRequest				body		PasswordUpdateRequest		true	"PasswordUpdateRequest"
// @Success			200									{object}	StatusResponse
// @Router			/users/update-password				[PUT]
func (s userAPI) PasswordUpdate(req *PasswordUpdateRequest) (res *StatusResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is empty")
	}
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.ConfirmNewPassword = strings.TrimSpace(req.ConfirmNewPassword)
	req.Mode = strings.ToLower(req.Mode)
	if req.Mode != modeChangePassword && req.Mode != modeSetPassword {
		return nil, fmt.Errorf("mode is invalid")
	}

	if len(req.NewPassword) < 8 {
		return nil, fmt.Errorf("new password min length is 8 chars")
	}
	if len(req.ConfirmNewPassword) < 8 {
		return nil, fmt.Errorf("new confirm password min length is 8 chars")
	}
	if req.NewPassword != req.ConfirmNewPassword {
		return nil, fmt.Errorf("new password and new confirm password doesn't match")
	}

	var user User
	result := s.db.Where("id = ?", req.UserID).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("finding user error: %s", result.Error.Error())
	}

	// Handle password set
	if req.Mode == modeSetPassword {
		if user.Password != "" {
			return nil, fmt.Errorf("password already has been set")
		}

	}

	if req.Mode == modeChangePassword {
		req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
		if req.ConfirmNewPassword == "" {
			return nil, fmt.Errorf("current password is required")
		}
		// Compare passwords
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
			return nil, fmt.Errorf("current password is incorrect")
		}
	}

	pwh, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("func: Signup, operation: bcrypt.GenerateFromPassword, err: %s", err.Error())
		return nil, err
	}
	pwhs := string(pwh)
	user.Password = pwhs

	result = s.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

// @Summary      	EmailUpdate
// @Description		Validates user id new email then will check if email is the same of exists then if not updates the email of the user.
// @Tags			Users
// @Produce			json
// @Param			Authorization						header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			EmailUpdateRequest					body		EmailUpdateRequest		true	"EmailUpdateRequest"
// @Success			200									{object}	StatusResponse
// @Router			/users/update-email					[PUT]
func (s userAPI) EmailUpdate(req *EmailUpdateRequest) (res *StatusResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is empty")
	}
	req.NewEmail = strings.TrimSpace(req.NewEmail)
	if req.NewEmail == "" {
		return nil, fmt.Errorf("new email is empty")
	}
	if !helpers.ValidEmail(req.NewEmail) {
		return nil, fmt.Errorf("new email not valid")
	}

	// find user by id
	var user User
	result := s.db.Where("id = ?", req.UserID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if req.NewEmail == user.Email {
		return nil, fmt.Errorf("the email is already the same")
	}

	var userEmail User
	_ = s.db.Where("email = ?", req.NewEmail).First(&userEmail)
	if userEmail.Email == req.NewEmail {
		return nil, fmt.Errorf("this email is already in use")
	}

	// update the email
	user.Email = req.NewEmail
	result = s.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

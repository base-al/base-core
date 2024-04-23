package users

import (
	helper "github.com/base-al/base-core/helpers"
	"github.com/base-al/base-core/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHTTPTransport interface {
	PasswordUpdate(c *fiber.Ctx) error
	EmailUpdate(c *fiber.Ctx) error
}

type userHttpTransport struct {
	userApi UserAPI
	logger  log.AllLogger
}

func NewUserHTTPTransport(uapi UserAPI, logger log.AllLogger) UserHTTPTransport {
	return &userHttpTransport{
		userApi: uapi,
		logger:  logger,
	}
}

func (s *userHttpTransport) PasswordUpdate(c *fiber.Ctx) error {
	req := &PasswordUpdateRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.PasswordUpdate.middleware.CtxUserID")
	}
	req.UserID = userId
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.PasswordUpdate.BodyParser")
	}
	resp, err := s.userApi.PasswordUpdate(req)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.PasswordUpdate.userApi.Update")
	}
	return c.JSON(resp)
}

func (s *userHttpTransport) EmailUpdate(c *fiber.Ctx) error {
	req := &EmailUpdateRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.EmailUpdate.middleware.CtxUserID")
	}
	req.UserID = userId
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.EmailUpdate.BodyParser")
	}
	resp, err := s.userApi.EmailUpdate(req)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.EmailUpdate.userApi.Update")
	}
	return c.JSON(resp)
}

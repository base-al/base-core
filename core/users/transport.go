package users

import (
	"base-core/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type HTTPTransport interface {
	Signup(c *fiber.Ctx) error
	SignupVerify(c *fiber.Ctx) error
}

type HttpTransport struct {
	userApi UserAPI
	logger  log.AllLogger
}

func NewHTTPTransport(uapi UserAPI, logger log.AllLogger) HTTPTransport {
	return &HttpTransport{
		userApi: uapi,
		logger:  logger,
	}
}

func (s *HttpTransport) Signup(c *fiber.Ctx) error {
	req := &SignupUserRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.Signup.BodyParser")
	}
	resp, err := s.userApi.Signup(req)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.Signup")
	}
	return c.JSON(resp)
}

func (s *HttpTransport) SignupVerify(c *fiber.Ctx) error {
	req := &SignupUserVerifyRequest{}
	req.Token = c.Params("token")
	resp, err := s.userApi.SignupVerify(req)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.SignupVerify")
	}
	return c.JSON(resp)
}

package auth

import (
	"base-core/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type HTTPTransport interface {
	Login(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
}

type httpTransport struct {
	authApi AuthAPI
	logger  log.AllLogger
}

func NewHTTPTransport(aapi AuthAPI, logger log.AllLogger) HTTPTransport {
	return &httpTransport{
		authApi: aapi,
		logger:  logger,
	}
}

func (s *httpTransport) Login(c *fiber.Ctx) error {
	req := &LoginRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.BodyParser")
	}
	resp, err := s.authApi.Login(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.Login")
	}
	return c.JSON(resp)
}

func (s *httpTransport) ForgotPassword(c *fiber.Ctx) error {
	req := &ForgotPasswordRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.BodyParser")
	}
	resp, err := s.authApi.ForgotPassword(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.ForgotPassword")
	}
	return c.JSON(resp)
}

func (s *httpTransport) ResetPassword(c *fiber.Ctx) error {
	req := &ResetPasswordRequest{}
	req.Token = c.Params("token")
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.BodyParser")
	}
	resp, err := s.authApi.ResetPassword(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.ResetPassword")
	}
	return c.JSON(resp)
}

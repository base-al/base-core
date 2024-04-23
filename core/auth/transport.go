package auth

import (
	"base/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthHTTPTransport interface {
	Signup(c *fiber.Ctx) error
	SignupVerify(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	OAuthLogin(c *fiber.Ctx) error
}

type authHttpTransport struct {
	authApi AuthAPI
	logger  log.AllLogger
}

func NewAuthHTTPTransport(aapi AuthAPI, logger log.AllLogger) AuthHTTPTransport {
	return &authHttpTransport{
		authApi: aapi,
		logger:  logger,
	}
}

func (s *authHttpTransport) Signup(c *fiber.Ctx) error {
	req := &SignupRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.Signup.BodyParser")
	}
	resp, err := s.authApi.Signup(req)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.Signup")
	}
	return c.JSON(resp)
}

func (s *authHttpTransport) SignupVerify(c *fiber.Ctx) error {
	req := &SignupVerifyRequest{}
	req.Token = c.Params("token")
	resp, err := s.authApi.SignupVerify(req)
	if err != nil {
		return helper.HTTPError(c, err, "UserHTTPTransport.SignupVerify")
	}
	return c.JSON(resp)
}

func (s *authHttpTransport) Login(c *fiber.Ctx) error {
	req := &LoginRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.Login.BodyParser")
	}
	resp, err := s.authApi.Login(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.Login")
	}
	return c.JSON(resp)
}

func (s *authHttpTransport) ForgotPassword(c *fiber.Ctx) error {
	req := &ForgotPasswordRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.ForgotPassword.BodyParser")
	}
	resp, err := s.authApi.ForgotPassword(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.ForgotPassword")
	}
	return c.JSON(resp)
}

func (s *authHttpTransport) ResetPassword(c *fiber.Ctx) error {
	req := &ResetPasswordRequest{}
	req.Token = c.Params("token")
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.ResetPassword.BodyParser")
	}
	resp, err := s.authApi.ResetPassword(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.ResetPassword")
	}
	return c.JSON(resp)
}

func (s *authHttpTransport) OAuthLogin(c *fiber.Ctx) error {
	req := &OAuthLoginRequest{}
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.OAuthLogin.BodyParser")
	}
	resp, err := s.authApi.OAuthLogin(req)
	if err != nil {
		return helper.HTTPError(c, err, "AuthHTTPTransport.authApi.OAuthLogin")
	}
	return c.JSON(resp)
}

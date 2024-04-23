package oauth

import (
	"base/helper"
	"base/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OAuthAccountHTTPTransport interface {
	GoogleSignUp(c *fiber.Ctx) error
	GoogleSignUpCallback(c *fiber.Ctx) error
	GoogleSignIn(c *fiber.Ctx) error
	GoogleSignInCallback(c *fiber.Ctx) error
	UserOAuthAccounts(c *fiber.Ctx) error
	GoogleAccountConnect(c *fiber.Ctx) error
	GoogleAccountConnectCallback(c *fiber.Ctx) error
}

type oAuthAccountHttpTransport struct {
	oAuthAccountApi OAuthAccountAPI
	logger          log.AllLogger
	uiAppUrl        string
}

func NewOAuthAccountHTTPTransport(oaapi OAuthAccountAPI, logger log.AllLogger, uiAppUrl string) OAuthAccountHTTPTransport {
	return &oAuthAccountHttpTransport{
		oAuthAccountApi: oaapi,
		logger:          logger,
		uiAppUrl:        uiAppUrl,
	}
}

func (s *oAuthAccountHttpTransport) GoogleSignUp(c *fiber.Ctx) error {
	resp, err := s.oAuthAccountApi.GoogleSignUp()
	if err != nil {
		return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.GoogleSignUp")
	}
	return c.JSON(resp)
}

func (s *oAuthAccountHttpTransport) GoogleSignUpCallback(c *fiber.Ctx) error {
	req := &OAuthCallbackRequest{}
	code := c.Query("code")
	req.Code = code
	resp, err := s.oAuthAccountApi.GoogleSignUpCallback(req)
	if err != nil {
		uiUrl := s.uiAppUrl + "/sign-up"
		uiUrl += "?googleoautherr=" + err.Error()
		return c.Redirect(uiUrl, fiber.StatusMovedPermanently)
		// return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.OAuthCallback")
	}
	uiUrl := s.uiAppUrl + "/login?provider=" + resp.Provider + "&t=" + resp.Token
	return c.Redirect(uiUrl)
}

func (s *oAuthAccountHttpTransport) GoogleSignIn(c *fiber.Ctx) error {
	resp, err := s.oAuthAccountApi.GoogleSignIn()
	if err != nil {
		return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.GoogleSignIn")
	}
	return c.JSON(resp)
}

func (s *oAuthAccountHttpTransport) GoogleSignInCallback(c *fiber.Ctx) error {
	req := &OAuthCallbackRequest{}
	code := c.Query("code")
	req.Code = code
	resp, err := s.oAuthAccountApi.GoogleSignInCallback(req)
	if err != nil {
		uiUrl := s.uiAppUrl + "/login"
		uiUrl += "?googleoautherr=" + err.Error()
		return c.Redirect(uiUrl, fiber.StatusMovedPermanently)
		// return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.OAuthCallback")
	}
	uiUrl := s.uiAppUrl + "/login?provider=" + resp.Provider + "&t=" + resp.Token
	return c.Redirect(uiUrl)
}

func (s *oAuthAccountHttpTransport) UserOAuthAccounts(c *fiber.Ctx) error {
	req := &FindUserOAuthAccount{}
	ctxUserId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.GoogleSignInCallback.middleware.CtxUserID")
	}
	req.UserID = ctxUserId

	resp, err := s.oAuthAccountApi.UserOAuthAccounts(req)
	if err != nil {
		helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.GoogleSignIn")
	}
	return c.JSON(resp)
}

func (s *oAuthAccountHttpTransport) GoogleAccountConnect(c *fiber.Ctx) error {
	req := &GoogleAccountConnectRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "UserSettingsHTTPTransport.CtxUserID")
	}
	req.BaseUserID = userId

	resp, err := s.oAuthAccountApi.GoogleAccountConnect(req)
	if err != nil {
		return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.GoogleAccountConnect")
	}
	return c.JSON(resp)
}

func (s *oAuthAccountHttpTransport) GoogleAccountConnectCallback(c *fiber.Ctx) error {
	req := &GoogleAccountConnectCallbackRequest{}
	code := c.Query("code")
	req.Code = code
	state := c.Query("state")
	req.State = state
	resp, err := s.oAuthAccountApi.GoogleAccountConnectCallback(req)
	if err != nil {
		return helper.HTTPError(c, err, "OAuthAccountHTTPTransport.oAuthAccountApi.GoogleAccountConnectCallback")
	}
	if resp.Status {
		return c.SendFile("./app/oauthaccounts/static/callback.html")
	}
	return c.JSON(fiber.Map{"error": err})
}

package google

import (
	"github.com/base-al/base-core/helpers"
	"github.com/base-al/base-core/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OAuthGoogleAccountHTTPTransport interface {
	GoogleSignUp(c *fiber.Ctx) error
	GoogleSignUpCallback(c *fiber.Ctx) error
	GoogleSignIn(c *fiber.Ctx) error
	GoogleSignInCallback(c *fiber.Ctx) error
	UserOAuthAccounts(c *fiber.Ctx) error
	GoogleAccountConnect(c *fiber.Ctx) error
	GoogleAccountConnectCallback(c *fiber.Ctx) error
}

type oAuthGoogleAccountHttpTransport struct {
	oAuthGoogleAccountApi OAuthGoogleAccountAPI
	logger                log.AllLogger
	uiAppUrl              string
}

func NewOAuthGoogleAccountHTTPTransport(oaapi OAuthGoogleAccountAPI, logger log.AllLogger, uiAppUrl string) OAuthGoogleAccountHTTPTransport {
	return &oAuthGoogleAccountHttpTransport{
		oAuthGoogleAccountApi: oaapi,
		logger:                logger,
		uiAppUrl:              uiAppUrl,
	}
}

func (s *oAuthGoogleAccountHttpTransport) GoogleSignUp(c *fiber.Ctx) error {
	resp, err := s.oAuthGoogleAccountApi.GoogleSignUp()
	if err != nil {
		return helpers.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.GoogleSignUp")
	}
	return c.JSON(resp)
}

func (s *oAuthGoogleAccountHttpTransport) GoogleSignUpCallback(c *fiber.Ctx) error {
	req := &OAuthCallbackRequest{}
	code := c.Query("code")
	req.Code = code
	resp, err := s.oAuthGoogleAccountApi.GoogleSignUpCallback(req)
	if err != nil {
		uiUrl := s.uiAppUrl + "/sign-up"
		uiUrl += "?googleoautherr=" + err.Error()
		return c.Redirect(uiUrl, fiber.StatusMovedPermanently)
		// return helper.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.OAuthCallback")
	}
	uiUrl := s.uiAppUrl + "/login?provider=" + resp.Provider + "&t=" + resp.Token
	return c.Redirect(uiUrl)
}

func (s *oAuthGoogleAccountHttpTransport) GoogleSignIn(c *fiber.Ctx) error {
	resp, err := s.oAuthGoogleAccountApi.GoogleSignIn()
	if err != nil {
		return helpers.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.GoogleSignIn")
	}
	return c.JSON(resp)
}

func (s *oAuthGoogleAccountHttpTransport) GoogleSignInCallback(c *fiber.Ctx) error {
	req := &OAuthCallbackRequest{}
	code := c.Query("code")
	req.Code = code
	resp, err := s.oAuthGoogleAccountApi.GoogleSignInCallback(req)
	if err != nil {
		uiUrl := s.uiAppUrl + "/login"
		uiUrl += "?googleoautherr=" + err.Error()
		return c.Redirect(uiUrl, fiber.StatusMovedPermanently)
		// return helper.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.OAuthCallback")
	}
	uiUrl := s.uiAppUrl + "/login?provider=" + resp.Provider + "&t=" + resp.Token
	return c.Redirect(uiUrl)
}

func (s *oAuthGoogleAccountHttpTransport) UserOAuthAccounts(c *fiber.Ctx) error {
	req := &FindUserOAuthAccount{}
	ctxUserId, err := middleware.CtxUserID(c)
	if err != nil {
		return helpers.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.GoogleSignInCallback.middleware.CtxUserID")
	}
	req.UserID = ctxUserId

	resp, err := s.oAuthGoogleAccountApi.UserOAuthAccounts(req)
	if err != nil {
		helpers.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.GoogleSignIn")
	}
	return c.JSON(resp)
}

func (s *oAuthGoogleAccountHttpTransport) GoogleAccountConnect(c *fiber.Ctx) error {
	req := &GoogleAccountConnectRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helpers.HTTPError(c, err, "UserSettingsHTTPTransport.CtxUserID")
	}
	req.BaseUserID = userId

	resp, err := s.oAuthGoogleAccountApi.GoogleAccountConnect(req)
	if err != nil {
		return helpers.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.GoogleAccountConnect")
	}
	return c.JSON(resp)
}

func (s *oAuthGoogleAccountHttpTransport) GoogleAccountConnectCallback(c *fiber.Ctx) error {
	req := &GoogleAccountConnectCallbackRequest{}
	code := c.Query("code")
	req.Code = code
	state := c.Query("state")
	req.State = state
	resp, err := s.oAuthGoogleAccountApi.GoogleAccountConnectCallback(req)
	if err != nil {
		return helpers.HTTPError(c, err, "OAuthGoogleAccountHTTPTransport.oAuthGoogleAccountApi.GoogleAccountConnectCallback")
	}
	if resp.Status {
		return c.SendFile("./app/oauthaccounts/static/callback.html")
	}
	return c.JSON(fiber.Map{"error": err})
}

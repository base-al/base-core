package oauth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, oAuthAccountHttpApi OAuthAccountHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	oAuthAccountRoutes := router.Group("oauth")
	oAuthAccountRoutes.Get("google/signup", oAuthAccountHttpApi.GoogleSignUp)
	oAuthAccountRoutes.Get("google/signup/callback", oAuthAccountHttpApi.GoogleSignUpCallback)
	oAuthAccountRoutes.Get("google", oAuthAccountHttpApi.GoogleSignIn)
	oAuthAccountRoutes.Get("google/callback", oAuthAccountHttpApi.GoogleSignInCallback)
	oAuthAccountRoutes.Get("accounts", authMiddleware, oAuthAccountHttpApi.UserOAuthAccounts)
	oAuthAccountRoutes.Get("google/connect", authMiddleware, oAuthAccountHttpApi.GoogleAccountConnect)
	oAuthAccountRoutes.Get("google/connect/callback", oAuthAccountHttpApi.GoogleAccountConnectCallback)
}

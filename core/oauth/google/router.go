package google

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, oAuthGoogleAccountHttpApi OAuthGoogleAccountHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	oAuthAccountRoutes := router.Group("oauth")
	oAuthAccountRoutes.Get("google/signup", oAuthGoogleAccountHttpApi.GoogleSignUp)
	oAuthAccountRoutes.Get("google/signup/callback", oAuthGoogleAccountHttpApi.GoogleSignUpCallback)
	oAuthAccountRoutes.Get("google", oAuthGoogleAccountHttpApi.GoogleSignIn)
	oAuthAccountRoutes.Get("google/callback", oAuthGoogleAccountHttpApi.GoogleSignInCallback)
	oAuthAccountRoutes.Get("accounts", authMiddleware, oAuthGoogleAccountHttpApi.UserOAuthAccounts)
	oAuthAccountRoutes.Get("google/connect", authMiddleware, oAuthGoogleAccountHttpApi.GoogleAccountConnect)
	oAuthAccountRoutes.Get("google/connect/callback", oAuthGoogleAccountHttpApi.GoogleAccountConnectCallback)
}

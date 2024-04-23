package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, authHttpApi AuthHTTPTransport) {
	authRoutes := router.Group("/auth")
	authRoutes.Post("/signup", authHttpApi.Signup)
	authRoutes.Put("/signup/verify/:token", authHttpApi.SignupVerify)
	authRoutes.Post("/login", authHttpApi.Login)
	authRoutes.Post("/forgot-password", authHttpApi.ForgotPassword)
	authRoutes.Post("/reset-password/:token", authHttpApi.ResetPassword)
	authRoutes.Post("/oauth", authHttpApi.OAuthLogin)
}

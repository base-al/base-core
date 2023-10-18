package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, authHttpApi HTTPTransport) {
	authRoutes := router.Group("/auth")
	authRoutes.Post("/login", authHttpApi.Login)
	authRoutes.Post("/forgot-password", authHttpApi.ForgotPassword)
	authRoutes.Put("/reset-password/:token", authHttpApi.ResetPassword)
}

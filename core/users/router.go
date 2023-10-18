package users

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, userHttpApi HTTPTransport) {
	userRoutes := router.Group("/users")
	userRoutes.Post("/signup", userHttpApi.Signup)
	userRoutes.Put("/signup/verify/:token", userHttpApi.SignupVerify)
}

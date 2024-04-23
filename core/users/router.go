package users

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, userHttpApi UserHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	// public
	userRoutes := router.Group("/users")
	// protected with auth
	userRoutes.Put("update-password", authMiddleware, userHttpApi.PasswordUpdate)
	userRoutes.Put("update-email", authMiddleware, userHttpApi.EmailUpdate)

}

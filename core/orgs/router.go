package orgs

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, authHttpApi HTTPTransport) {
	userRoutes := router.Group("/orgs")
	userRoutes.Post("/", authHttpApi.Add)
	userRoutes.Get("/me", authHttpApi.FindMyOrgs)
}

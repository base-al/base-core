// app/pages/router.go

package pages

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterPageRoutes sets up the routing for page management.
func RegisterPageRoutes(router fiber.Router, transport *PageHTTPTransport) {
	// Here you can also add middleware specific to these routes
	//index
	router.Get("/", transport.Index)
	router.Post("/", transport.Create)      // Handles POST requests to create a new page
	router.Get("/:id", transport.Read)      // Handles GET requests to read a page by ID
	router.Put("/:id", transport.Update)    // Handles PUT requests to update a page by ID
	router.Delete("/:id", transport.Delete) // Handles DELETE requests to delete a page by ID
}

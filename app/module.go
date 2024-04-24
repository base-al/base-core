// app/module.go

package app

import (
	"github.com/base-al/base-core/app/pages"
	"github.com/gofiber/fiber/v2"
)

func RegisterAllModules(app *fiber.App) {
	// Register everything for the Pages module
	// Instantiate the service and transport layer
	pageService := pages.NewSimplePageService()              // Instantiate the service
	pageTransport := pages.NewPageHTTPTransport(pageService) // Instantiate the transport
	// Register the routes
	pagesRouter := app.Group("/pages")
	pages.RegisterPageRoutes(pagesRouter, pageTransport)

}

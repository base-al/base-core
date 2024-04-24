package chapter

import (
    "github.com/gofiber/fiber/v2"
)

// RegisterChapterRoutes sets up the routing for chapter management.
func RegisterChapterRoutes(router fiber.Router, transport *ChapterHTTPTransport) {
    router.Get("/", transport.Index)       // Handles GET requests to list chapters
    router.Post("/", transport.Create)     // Handles POST requests to create a new chapter
    router.Get("/:id", transport.Read)     // Handles GET requests to read a chapter by ID
    router.Put("/:id", transport.Update)   // Handles PUT requests to update a chapter by ID
    router.Delete("/:id", transport.Delete) // Handles DELETE requests to delete a chapter by ID
}

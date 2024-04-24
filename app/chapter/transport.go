package chapter

import (
	"github.com/gofiber/fiber/v2"
)

type ChapterHTTPTransport struct {
	service ChapterService
}

func NewChapterHTTPTransport(service ChapterService) *ChapterHTTPTransport {
	return &ChapterHTTPTransport{service: service}
}

func (p *ChapterHTTPTransport) RegisterRoutes(router fiber.Router) {
	router.Get("/", p.Index)        // List all chapters
	router.Post("/", p.Create)      // Create a new chapter
	router.Get("/:id", p.Read)      // Get a chapter by ID
	router.Put("/:id", p.Update)    // Update a chapter by ID
	router.Delete("/:id", p.Delete) // Delete a chapter by ID
}

func (p *ChapterHTTPTransport) Index(c *fiber.Ctx) error {
	resp, err := p.service.Index()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list chapters"})
	}
	return c.JSON(resp)
}

func (p *ChapterHTTPTransport) Create(c *fiber.Ctx) error {
	var req CreateChapterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}
	resp, err := p.service.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create chapter"})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (p *ChapterHTTPTransport) Read(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	req := ReadChapterRequest{ID: id}
	resp, err := p.service.Read(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "chapter not found"})
	}
	return c.JSON(resp)
}

func (p *ChapterHTTPTransport) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	var req UpdateChapterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}
	req.ID = id
	resp, err := p.service.Update(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update chapter"})
	}
	return c.JSON(resp)
}

func (p *ChapterHTTPTransport) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	req := DeleteChapterRequest{ID: id}
	resp, err := p.service.Delete(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete chapter"})
	}
	return c.JSON(resp)
}

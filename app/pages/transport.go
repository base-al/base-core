package pages

import (
	"github.com/gofiber/fiber/v2"
)

type PageHTTPTransport struct {
	service PageService
}

func NewPageHTTPTransport(service PageService) *PageHTTPTransport {
	return &PageHTTPTransport{service: service}
}

func (p *PageHTTPTransport) RegisterRoutes(router fiber.Router) {
	router.Get("/", p.Index)        // List all pages
	router.Post("/", p.Create)      // Create a new page
	router.Get("/:id", p.Read)      // Get a page by ID
	router.Put("/:id", p.Update)    // Update a page by ID
	router.Delete("/:id", p.Delete) // Delete a page by ID
}

func (p *PageHTTPTransport) Index(c *fiber.Ctx) error {
	resp, err := p.service.Index()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list pages"})
	}
	return c.JSON(resp)
}

func (p *PageHTTPTransport) Create(c *fiber.Ctx) error {
	var req CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}
	resp, err := p.service.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create page"})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (p *PageHTTPTransport) Read(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	req := ReadPageRequest{ID: id}
	resp, err := p.service.Read(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "page not found"})
	}
	return c.JSON(resp)
}

func (p *PageHTTPTransport) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	var req UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}
	req.ID = id
	resp, err := p.service.Update(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update page"})
	}
	return c.JSON(resp)
}

func (p *PageHTTPTransport) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	req := DeletePageRequest{ID: id}
	resp, err := p.service.Delete(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete page"})
	}
	return c.JSON(resp)
}

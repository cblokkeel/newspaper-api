package categories

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type CategoriesHandler struct {
	svc *CategoriesService
}

func NewCategoriesHandler(svc *CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{
		svc,
	}
}

func (h *CategoriesHandler) Mount(fiber fiber.Router) {
	group := fiber.Group("/categories")
	group.Get("/", h.GetLocalisedCategories)
}

func (h *CategoriesHandler) GetLocalisedCategories(c *fiber.Ctx) error {
	locale := c.Query("l", "fr") // TODO change for en
	if locale != "fr" && locale != "en" {
		locale = "fr"
	}
	categories, err := h.svc.getCategories(c.Context(), locale)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "something went wrong"})
	}
	return c.JSON(categories)
}

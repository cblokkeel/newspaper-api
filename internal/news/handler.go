package news

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct {
	svc *NewsService
}

func NewNewsHandler(svc *NewsService) *NewsHandler {
	return &NewsHandler{
		svc,
	}
}

func (h *NewsHandler) Mount(fiber fiber.Router) {
	group := fiber.Group("/news")
	group.Get("/", h.HandleGetNews)
	group.Get("/sources", h.HandleGetSources)
}

func (h *NewsHandler) HandleGetNews(c *fiber.Ctx) error {
	news, err := h.svc.getNews(c.Context())
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(news)
}

func (h *NewsHandler) HandleGetSources(c *fiber.Ctx) error {
	// // sources, err := h.svc.getSources(c.Context())
	// if err != nil {
	// 	return c.SendStatus(http.StatusInternalServerError)
	// }
	//
	return c.SendString("todo")
}

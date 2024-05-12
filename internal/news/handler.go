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

func (h *NewsHandler) HandleGetNews(c *fiber.Ctx) error {
	news, err := h.svc.GetNews(c.Context())
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(news)
}

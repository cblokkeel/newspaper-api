package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/cblokkeel/newspaper/internal/deps"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Everything ok")
	})

	{
		newsAPI := api.Group("/news")
		newsHandler := deps.GetNewsHandler()
		newsAPI.Get("/", newsHandler.HandleGetNews)
		newsAPI.Get("/sources", newsHandler.HandleGetSources)
	}

    {
        categoriesAPI := api.Group("/categories")
        categoriesHandler := deps.GetCategoriesHandler()
    categoriesAPI.Get("/", categoriesHandler.GetLocalisedCategories)
    }

	return app
}

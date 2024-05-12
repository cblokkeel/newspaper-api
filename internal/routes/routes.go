package routes

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/cblokkeel/newspaper/internal/news"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := rdb.Set(context.Background(), "foo", "bar", 0).Err(); err != nil {
		panic(err)
	}

    newsSvc := news.NewNewsService(rdb)
    newsHandler := news.NewNewsHandler(newsSvc)

	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Everything ok")
	})

	newsAPI := api.Group("/news")
	newsAPI.Get("/", newsHandler.HandleGetNews)

	return app
}

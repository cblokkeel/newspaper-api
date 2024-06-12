package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cblokkeel/newspaper/internal/categories"
	"github.com/cblokkeel/newspaper/internal/db/mongo"
	redisdb "github.com/cblokkeel/newspaper/internal/db/redis"
	weaviatedb "github.com/cblokkeel/newspaper/internal/db/weaviate"
	"github.com/cblokkeel/newspaper/internal/httpsrv"
	"github.com/cblokkeel/newspaper/internal/news"
	rsscron "github.com/cblokkeel/newspaper/internal/rss_cron"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"go.uber.org/fx"
)

func http(app *fiber.App) {
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT"))); err != nil {
		panic(err)
	}
}

func newFiber(handlers []httpsrv.Handler) *fiber.App {
	app := fiber.New()
	for _, handler := range handlers {
		handler.Mount(app)
	}
	return app
}

// Temp
func newCron(mongo *mongo.MongoDB, weaviate *weaviatedb.WeaviateDB, app *fiber.App) {
	rssCron := rsscron.NewRSSCron(mongo, weaviate)
	c := cron.New()
	err := c.AddFunc("@every 1h", func() {
		rssCron.Start(context.Background())
	})
	if err != nil {
		log.Fatalf("Error scheduling rss cron job: %+v", err)
	}
	c.Start()
	defer c.Stop()
	app.Get("/cron/rss", func(fc *fiber.Ctx) error {
		fmt.Println("Starting")
		rssCron.Start(fc.Context())
		return fc.SendString("Job done")
	})
}

func AsHandler(h any) any {
	return fx.Annotate(
		h,
		fx.As(new(httpsrv.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}

func main() {
	fx.New(
		fx.Provide(
			// Misc
			redisdb.NewRedisDB,
			mongo.NewMongoDB,
			// Http
			fx.Annotate(
				newFiber,
				fx.ParamTags(`group:"handlers"`),
			),
			// Svc
			news.NewNewsService,
			categories.NewCategoriesService,
			// Handlers
			AsHandler(news.NewNewsHandler),
			AsHandler(categories.NewCategoriesHandler),
		),
		fx.Invoke(func(app *fiber.App) {
			http(app)
		}),
	).Run()
}

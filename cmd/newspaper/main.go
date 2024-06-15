package main

import (
	"fmt"
	"os"

	"github.com/cblokkeel/newspaper/internal/categories"
	"github.com/cblokkeel/newspaper/internal/db/mongo"
	redisdb "github.com/cblokkeel/newspaper/internal/db/redis"
	weaviatedb "github.com/cblokkeel/newspaper/internal/db/weaviate"
	"github.com/cblokkeel/newspaper/internal/httpsrv"
	"github.com/cblokkeel/newspaper/internal/news"
	"github.com/cblokkeel/newspaper/internal/workers"
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

func apiRouter(app *fiber.App) fiber.Router {
	return app.Group("/api")
}

func AsHandler(h any) any {
	return fx.Annotate(
		h,
		fx.As(new(httpsrv.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}

func AsWorker(w any) any {
	return fx.Annotate(
		w,
		fx.As(new(workers.Worker)),
		fx.ResultTags(`group:"workers"`),
	)
}

func AsWeaviateSchema(s any) any {
	return fx.Annotate(
		s,
		fx.As(new(weaviatedb.WeaviateSchema)),
		fx.ResultTags(`group:"w_schemas"`),
	)
}

func main() {
	fx.New(
		fx.Provide(
			// Misc
			redisdb.NewRedisDB,
			mongo.NewMongoDB,
			AsWeaviateSchema(weaviatedb.GetArticleSchema),
			fx.Annotate(
				weaviatedb.NewWeaviateDB,
				fx.ParamTags(`group:"w_schemas"`),
			),
			// Http
			fx.Annotate(
				newFiber,
				fx.ParamTags(`group:"handlers"`),
			),
			apiRouter,
			// Svc
			news.NewNewsService,
			categories.NewCategoriesService,
			// Handlers
			AsHandler(news.NewNewsHandler),
			AsHandler(categories.NewCategoriesHandler),
			// Workers
			AsWorker(workers.NewRssWorker),
			fx.Annotate(
				workers.NewWorkerManager,
				fx.ParamTags(`group:"workers"`),
			),
			workers.StartWorkers,
		),
		fx.Invoke(func(*cron.Cron) {}),
		fx.Invoke(func(app *fiber.App) {
			go http(app)
		}),
	).Run()
}

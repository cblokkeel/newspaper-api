package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"

	"github.com/cblokkeel/newspaper/internal/db"
	"github.com/cblokkeel/newspaper/internal/routes"
	rsscron "github.com/cblokkeel/newspaper/internal/rss_cron"
)

func main() {
	app := routes.NewRouter()
	mongo, err := db.NewMongoDB(os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
	defer mongo.Disconnect(context.Background())
	if err != nil {
		panic(fmt.Errorf("Couldn't connect to mongo"))
	}

	rssCron := rsscron.NewRSSCron(mongo)
	c := cron.New()
	_, err = c.AddFunc("@every 1h", func() {
		rssCron.Start(context.Background())
	})
	if err != nil {
		log.Fatalf("Error scheduling rss cron job: %+v", err)
	}

    c.Start()
    defer c.Stop()

    app.Get("/cron/rss", func (fc *fiber.Ctx) error {
        fmt.Println("Starting")
        rssCron.Start(fc.Context())
        return fc.SendString("Job done")
    })

	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT"))); err != nil {
		panic(err)
	}
}

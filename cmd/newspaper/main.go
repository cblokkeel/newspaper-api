package main

import (
	"fmt"
	"os"

	"github.com/cblokkeel/newspaper/internal/routes"
	_ "github.com/joho/godotenv/autoload"
)


func main() {
    app := routes.NewRouter()

    if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT"))); err != nil {
        panic(err)
    } 
}

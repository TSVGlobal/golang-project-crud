package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := App{}
	if err := app.Initialize(); err != nil {
		log.Fatalf("could not initialize app: %v", err)
	}
	app.handleRoutes()
	app.Run(os.Getenv("APP_ADDR"))
}

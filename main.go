package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	//"github.com/joho/godotenv"
	"github.com/rabbice/url_shortener/routes"
)

func main() {
	/*err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	app := fiber.New()
	app.Use(logger.New())

	routes.InitRoutes(app)

	log.Fatal(app.Listen(os.Getenv("PORT")))
}

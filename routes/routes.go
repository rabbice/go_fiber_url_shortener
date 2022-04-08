package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rabbice/url_shortener/handlers"
)

func InitRoutes(app *fiber.App) {
	app.Get("/:url", handlers.ResolveURL)
	app.Post("/v1/api", handlers.ShortenURL)
}

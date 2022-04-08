package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbice/url_shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.RedisConnection(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error in DB connection"})
	}

	count := database.RedisConnection(1)
	defer count.Close()
	_ = count.Incr(database.Ctx, "counter")
	return c.Redirect(value, 301)
}

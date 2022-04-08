package handlers

import (
	//"log"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	//"github.com/joho/godotenv"
	"github.com/rabbice/url_shortener/database"
	"github.com/rabbice/url_shortener/helpers"
	"github.com/rabbice/url_shortener/models"
)

func ShortenURL(c *fiber.Ctx) error {
	/*err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	body := new(models.Request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// rate limiting
	r1 := database.RedisConnection(1)
	defer r1.Close()
	val, err := r1.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r1.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r1.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceeded", "rate_limit": limit / time.Nanosecond / time.Minute})
		}
	}
	// also check if input is URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input. Input should be URL"})
	}
	// check domain
	if !helpers.RemoveDomainError(body.URL) {
		c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "You can't access this"})
	}
	// enforce TLS
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.RedisConnection(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "The URL is already used"})
	}
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "can't connect to server"})
	}

	resp := models.Response{
		URL:         body.URL,
		CustomShort: "",
		Expiry:      body.Expiry,
		RateLimit:   30,
		RateRemain:  10,
	}

	r1.Decr(database.Ctx, c.IP())

	val, _ = r1.Get(database.Ctx, c.IP()).Result()
	resp.RateRemain, _ = strconv.Atoi(val)
	ttl, _ := r1.TTL(database.Ctx, c.IP()).Result()
	resp.RateLimit = ttl / time.Nanosecond / time.Minute
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}

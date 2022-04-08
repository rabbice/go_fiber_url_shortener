package database

import (
	"context"
	//"log"
	"os"

	"github.com/go-redis/redis/v8"
	//"github.com/joho/godotenv"
)

var Ctx = context.Background()

func RedisConnection(dbno int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       dbno,
	})
	return rdb

}

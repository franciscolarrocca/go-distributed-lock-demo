package tools

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisSync() *redsync.Redsync {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Panicf("Error: Failed to retrieve REDIS_DB variable: %v", err)
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	options := &redis.Options{
		Addr:     redisAddr,
		DB:       redisDB,
		Password: redisPassword,
	}

	client := redis.NewClient(options)
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return rs
}

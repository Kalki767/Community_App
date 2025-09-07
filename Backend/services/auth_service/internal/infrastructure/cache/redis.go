package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		addr = "localhost:6379" // default
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Set password if needed
		DB:       0,  // default DB
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
	return rdb
}

package db

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func CreateRedisClient() {
	redisURL := os.Getenv("REDIS")

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("failed to parse redis url: ", err.Error())
	}

	client = redis.NewClient(opt)
}

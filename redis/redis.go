package db

import (
	"context"
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

func CreatePubSubClient(ctx context.Context, id string) *redis.PubSub {
	return client.Subscribe(ctx, id)
}

func PublishToRedis(ctx context.Context, id, data string) {
	err := client.Publish(ctx, id, data).Err()
	if err != nil {
		log.Println("Failed to publish to redis: ", err.Error())
	}
}

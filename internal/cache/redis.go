package cache

import (
	"context"
	"os"

	"starter/pkg/logger"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	log := logger.GetLogger()
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379/0"
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Redis URL")
	}

	RedisClient = redis.NewClient(opts)

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	log.Info().Msg("Redis connection established")
}

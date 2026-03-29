package config

import (
	"context"
	"fmt"
	"os"

	"github.com/NhatHaoDev3324/zizone-be/factory"
	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
	Ctx   = context.Background()
)

func ConnectRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		factory.LogError("Failed to connect to Redis: " + err.Error())
	}

	factory.LogSuccess("Connected to Redis successfully!")
	return Redis
}

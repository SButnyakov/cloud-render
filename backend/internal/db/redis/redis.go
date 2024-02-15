package redis

import (
	"cloud-render/internal/lib/config"
	"context"

	"github.com/redis/go-redis/v9"
)

func New(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
		DB:   0,
	})
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}

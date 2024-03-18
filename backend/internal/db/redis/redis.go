package redis

import (
	"cloud-render/internal/lib/config"
	"context"
	"errors"
	"time"

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

func Clear(client *redis.Client, queueName string) {
	var err error
	for {
		_, err = client.BLPop(context.Background(), time.Second, queueName).Result()
		if errors.Is(err, redis.Nil) {
			return
		}
	}
}

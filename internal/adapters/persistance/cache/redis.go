package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{
		client: client,
	}
}

func (r *RedisAdapter) Set(key string, value interface{}) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, 24*time.Hour).Err()
}

func (r *RedisAdapter) Get(key string) (interface{}, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

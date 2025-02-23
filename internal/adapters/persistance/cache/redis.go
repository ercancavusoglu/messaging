package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

type RedisAdapter struct {
	client RedisClient
}

func NewRedisAdapter(client RedisClient) *RedisAdapter {
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

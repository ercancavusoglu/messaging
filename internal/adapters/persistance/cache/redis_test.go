package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisAdapter_Integration(t *testing.T) {
	// Redis bağlantısı gerektiği için bu testi skip edelim
	t.Skip("Skipping integration test")

	// Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer client.Close()

	// Bağlantıyı test et
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	assert.NoError(t, err)

	// Redis adapter
	adapter := NewRedisAdapter(client)

	// Test verileri
	key := "test_key"
	value := "test_value"

	// Set
	err = adapter.Set(key, value)
	assert.NoError(t, err)

	// Get
	result, err := adapter.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Olmayan key
	_, err = adapter.Get("nonexistent_key")
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

type mockRedisClient struct {
	redis.Client
	mockSet func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	mockGet func(ctx context.Context, key string) *redis.StringCmd
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.mockSet(ctx, key, value, expiration)
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.mockGet(ctx, key)
}

func TestRedisAdapter_Set(t *testing.T) {
	// Test verileri
	key := "test_key"
	value := "test_value"

	// Mock client
	mockClient := &mockRedisClient{
		mockSet: func(ctx context.Context, k string, v interface{}, expiration time.Duration) *redis.StatusCmd {
			assert.Equal(t, key, k)
			assert.Equal(t, value, v)
			return redis.NewStatusCmd(ctx)
		},
	}

	// Redis adapter
	adapter := NewRedisAdapter(mockClient)

	// Set
	err := adapter.Set(key, value)
	assert.NoError(t, err)
}

func TestRedisAdapter_Get(t *testing.T) {
	// Test verileri
	key := "test_key"
	value := "test_value"

	// Mock client
	mockClient := &mockRedisClient{
		mockGet: func(ctx context.Context, k string) *redis.StringCmd {
			assert.Equal(t, key, k)
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal(value)
			return cmd
		},
	}

	// Redis adapter
	adapter := NewRedisAdapter(mockClient)

	// Get
	result, err := adapter.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestRedisAdapter_Get_NotFound(t *testing.T) {
	// Test verileri
	key := "nonexistent_key"

	// Mock client
	mockClient := &mockRedisClient{
		mockGet: func(ctx context.Context, k string) *redis.StringCmd {
			assert.Equal(t, key, k)
			cmd := redis.NewStringCmd(ctx)
			cmd.SetErr(redis.Nil)
			return cmd
		},
	}

	// Redis adapter
	adapter := NewRedisAdapter(mockClient)

	// Get
	_, err := adapter.Get(key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

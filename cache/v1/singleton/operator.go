package singleton

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// Operator is an implement that operates singleton redis
type Operator struct {
	SingletonClient *redis.Client
}

// Set Redis `SET key value [expiration]` command.
func (s Operator) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.SingletonClient.Set(ctx, key, value, expiration).Err()
}

// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
func (s Operator) Get(ctx context.Context, key string) (string, error) {
	return s.SingletonClient.Get(ctx, key).Result()
}

// HGet Redis `HGET key field` command.
func (s Operator) HGet(ctx context.Context, key, field string) (string, error) {
	return s.SingletonClient.HGet(ctx, key, field).Result()
}

// HGetAll Redis `HGETALL key` command.
func (s Operator) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return s.SingletonClient.HGetAll(ctx, key).Result()
}

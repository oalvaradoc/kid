package cluster

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// Operator is an implement that operates cluster redis
type Operator struct {
	ClusterClient *redis.ClusterClient
}

// Set Redis `SET key value [expiration]` command.
func (c Operator) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.ClusterClient.Set(ctx, key, value, expiration).Err()
}

// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
func (c Operator) Get(ctx context.Context, key string) (string, error) {
	return c.ClusterClient.Get(ctx, key).Result()
}

// HGet Redis `HGET key field` command.
func (c Operator) HGet(ctx context.Context, key, field string) (string, error) {
	return c.ClusterClient.HGet(ctx, key, field).Result()
}

// HGetAll Redis `HGETALL key` command.
func (c Operator) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.ClusterClient.HGetAll(ctx, key).Result()
}

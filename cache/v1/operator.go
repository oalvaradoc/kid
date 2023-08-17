package cache

import (
	"context"
	"time"
)

var AddressingCacheOperator Operator

// Operator is an interface that defines all the cache related functions
type Operator interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}

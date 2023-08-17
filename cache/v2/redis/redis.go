package redis

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"time"
)

// NewRedisConnectionPoolCache creates a Redis connection pool cache
func NewRedisConnectionPoolCache() *_cache {
	var redisConnectionPoolCache = &_cache{
		cache: make(map[string]redis.UniversalClient),
	}

	return redisConnectionPoolCache
}

type _cache struct {
	sync.RWMutex
	cache map[string]redis.UniversalClient
}

func (e *_cache) Get(aliasName string) (eg interface{}, ok bool) {
	e.RLock()
	defer e.RUnlock()
	if len(e.cache) > 0 {
		eg, ok = e.cache[aliasName]
	}

	return
}

func (e *_cache) set(aliasName string, redisClient redis.UniversalClient) {
	e.Lock()
	defer e.Unlock()
	if nil == e.cache {
		e.cache = make(map[string]redis.UniversalClient)
	}

	e.cache[aliasName] = redisClient
}

func (e *_cache) Delete(aliasName string) {
	eg, ok := e.Get(aliasName)
	if ok {
		// close Cache
		eg.(redis.UniversalClient).Close()
	}
	if nil != e.cache {
		delete(e.cache, aliasName)
	}
}

func (e *_cache) InitCache(aliasName string, cacheConfig *config.Cache) (err *errors.Error) {
	poolSize := constant.DefaultRedisPoolSize
	if cacheConfig.Pool.PoolSize > 0 {
		poolSize = cacheConfig.Pool.PoolSize
	}

	universalOptions := &redis.UniversalOptions{
		Addrs:    strings.Split(cacheConfig.Addr, ","),
		Password: cacheConfig.Password,
		PoolSize: poolSize,
	}

	if cacheConfig.Pool.IdleTimeoutSeconds > 0 {
		universalOptions.IdleTimeout = time.Second * time.Duration(cacheConfig.Pool.IdleTimeoutSeconds)
	}

	if cacheConfig.Pool.PoolTimeoutSeconds > 0 {
		universalOptions.PoolTimeout = time.Second * time.Duration(cacheConfig.Pool.PoolTimeoutSeconds)
	}

	if cacheConfig.Pool.MinIdleConns > 0 {
		universalOptions.MinIdleConns = cacheConfig.Pool.MinIdleConns
	}

	if cacheConfig.Pool.MaxConnAgeSeconds > 0 {
		universalOptions.MaxConnAge = time.Second * time.Duration(cacheConfig.Pool.MaxConnAgeSeconds)
	}

	if cacheConfig.Pool.DialTimeoutSeconds > 0 {
		universalOptions.DialTimeout = time.Second * time.Duration(cacheConfig.Pool.DialTimeoutSeconds)
	}

	if cacheConfig.Pool.ReadTimeoutSeconds > 0 {
		universalOptions.ReadTimeout = time.Second * time.Duration(cacheConfig.Pool.ReadTimeoutSeconds)
	}

	if cacheConfig.Pool.WriteTimeoutSeconds > 0 {
		universalOptions.WriteTimeout = time.Second * time.Duration(cacheConfig.Pool.WriteTimeoutSeconds)
	}

	redisClient := redis.NewUniversalClient(universalOptions)
	_, ge := redisClient.Ping(context.Background()).Result()
	if ge != nil {
		return errors.Wrap(constant.SystemInternalError, ge, 0)
	}

	e.set(aliasName, redisClient)

	log.Infosf("Successfully to init cache[%++v] for redis", *cacheConfig)
	return nil
}

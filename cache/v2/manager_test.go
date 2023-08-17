package cache

import (
	"context"
	"git.multiverse.io/eventkit/kit/handler/config"
	"testing"
)

func TestRedisManagement(t *testing.T) {
	t.Skip()
	cacheConfigs := map[string]config.Cache{
		"su0001-cache1": config.Cache{
			Name:     "",
			Type:     "redis",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:6379",
			Password: "",
			Pool: struct {
				PoolSize           int `json:"poolSize"`
				MinIdleConns       int `json:"minIdleConns"`
				MaxConnAgeSeconds  int `json:"maxConnAgeSeconds"`
				PoolTimeoutSeconds int `json:"poolTimeoutSeconds"`
				IdleTimeoutSeconds int `json:"idleTimeoutSeconds"`
				DialTimeoutSeconds int `json:"dialTimeoutSeconds"`
				ReadTimeoutSeconds int `json:"readTimeoutSeconds"`
				WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
			}{
				PoolSize: 20,
			},
		}, "su0002-cache1": config.Cache{
			Name:     "",
			Type:     "redis",
			Su:       "su0002",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:6379",
			Password: "",
			Pool: struct {
				PoolSize           int `json:"poolSize"`
				MinIdleConns       int `json:"minIdleConns"`
				MaxConnAgeSeconds  int `json:"maxConnAgeSeconds"`
				PoolTimeoutSeconds int `json:"poolTimeoutSeconds"`
				IdleTimeoutSeconds int `json:"idleTimeoutSeconds"`
				DialTimeoutSeconds int `json:"dialTimeoutSeconds"`
				ReadTimeoutSeconds int `json:"readTimeoutSeconds"`
				WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
			}{
				PoolSize: 20,
			},
		},
	}

	if gerr := InitCacheManagerForRedis(cacheConfigs); nil != gerr {
		t.Logf("Failed to init cache manager for redis, error:%++v", gerr)
	}

	// select from cache1
	redisClient, err := GetRedisClient("su0001", "DAS000011")
	if nil != err {
		t.Errorf("[cache1]Failed to get redis client, error:%++v", err)
		return
	}

	redisClient.HSet(context.Background(), "test-key-from-tester1", "sub-test-key-from-tester1", "test-value-from-tester1")
	v := redisClient.HGet(context.Background(), "test-key-from-tester1", "sub-test-key-from-tester1")
	t.Logf("[cache2]The result of Test record:%++v", v.String())

	// select from cache2
	redisClient, err = GetRedisClient("su0002", "DAS000011")
	if nil != err {
		t.Errorf("[cache2]Failed to get redis client, error:%++v", err)
		return
	}

	redisClient.HSet(context.Background(), "test-key-from-tester2", "sub-test-key-from-tester2", "test-value-from-tester2")
	v = redisClient.HGet(context.Background(), "test-key-from-tester2", "sub-test-key-from-tester2")
	t.Logf("[cache2]The result of Test record:%++v", v.String())

	cacheConfigs = map[string]config.Cache{
		"su0001-cache1": config.Cache{
			Name:     "",
			Type:     "redis",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:6379",
			Password: "",
			Pool: struct {
				PoolSize           int `json:"poolSize"`
				MinIdleConns       int `json:"minIdleConns"`
				MaxConnAgeSeconds  int `json:"maxConnAgeSeconds"`
				PoolTimeoutSeconds int `json:"poolTimeoutSeconds"`
				IdleTimeoutSeconds int `json:"idleTimeoutSeconds"`
				DialTimeoutSeconds int `json:"dialTimeoutSeconds"`
				ReadTimeoutSeconds int `json:"readTimeoutSeconds"`
				WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
			}{
				PoolSize: 10,
			},
		}, "su0001-cache2": config.Cache{
			Name:     "",
			Type:     "redis",
			Su:       "su0001",
			Topics:   []string{"DAS00001", "DAS00002"},
			Default:  false,
			Addr:     "127.0.0.1:6379",
			Password: "",
			Pool: struct {
				PoolSize           int `json:"poolSize"`
				MinIdleConns       int `json:"minIdleConns"`
				MaxConnAgeSeconds  int `json:"maxConnAgeSeconds"`
				PoolTimeoutSeconds int `json:"poolTimeoutSeconds"`
				IdleTimeoutSeconds int `json:"idleTimeoutSeconds"`
				DialTimeoutSeconds int `json:"dialTimeoutSeconds"`
				ReadTimeoutSeconds int `json:"readTimeoutSeconds"`
				WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
			}{
				PoolSize: 20,
			},
		}, "su0002-cache1": config.Cache{
			Name:     "",
			Type:     "redis",
			Su:       "su0002",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:6379",
			Password: "",
			Pool: struct {
				PoolSize           int `json:"poolSize"`
				MinIdleConns       int `json:"minIdleConns"`
				MaxConnAgeSeconds  int `json:"maxConnAgeSeconds"`
				PoolTimeoutSeconds int `json:"poolTimeoutSeconds"`
				IdleTimeoutSeconds int `json:"idleTimeoutSeconds"`
				DialTimeoutSeconds int `json:"dialTimeoutSeconds"`
				ReadTimeoutSeconds int `json:"readTimeoutSeconds"`
				WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
			}{
				PoolSize: 20,
			},
		},
	}

	if gerr := Rotate(cacheConfigs); nil != gerr {
		t.Errorf("Failed to rotate the cache configuration, error:%++v", gerr)
	}

	redisClient, err = GetRedisClient("su0001", "DAS000011")
	if nil != err {
		t.Errorf("[cache1]Failed to get redis client, error:%++v", err)
		return
	}

	t.Logf("1.redisClient:%++p", redisClient)

	redisClient2, err := GetRedisClient("su0001", "DAS00001")
	if nil != err {
		t.Errorf("[cache1]Failed to get redis client, error:%++v", err)
		return
	}
	t.Logf("2.redisClient:%++p", redisClient2)

	redisClient3, err := GetRedisClient("su0001", "DAS00003")
	if nil != err {
		t.Errorf("[cache1]Failed to get redis client, error:%++v", err)
		return
	}
	t.Logf("3.redisClient:%++p", redisClient3)

	if redisClient != redisClient3 {
		t.Errorf("Except same client")
	}

	if redisClient2 == redisClient3 {
		t.Errorf("Except different client")
	}

	redisClient, err = GetRedisClient("su0002", "DAS000011")
	if nil != err {
		t.Errorf("[cache2]Failed to get redis client, error:%++v", err)
		return
	}
	t.Logf("4.redisClient:%++p", redisClient)
}

package repository

import (
	"context"
	"git.multiverse.io/eventkit/kit/cache/v1"
	"git.multiverse.io/eventkit/kit/cache/v1/cluster"
	"git.multiverse.io/eventkit/kit/cache/v1/singleton"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model/glsdef"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

// CacheRepository is an implement that adapter singleton or cluster redis operator
type CacheRepository struct {
	cacheOperator cache.Operator
}

func (c CacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.cacheOperator.Set(ctx, key, value, expiration)
}

func (c CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return c.cacheOperator.Get(ctx, key)
}

func (c CacheRepository) HGet(ctx context.Context, key, field string) (string, error) {
	return c.cacheOperator.HGet(ctx, key, field)
}

func (c CacheRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.cacheOperator.HGetAll(ctx, key)
}

// newSingletonRepository creates a new singleton redis operator via redis.Options
func newSingletonRepository(options *redis.Options) cache.Operator {
	return &CacheRepository{singleton.Operator{SingletonClient: redis.NewClient(options)}}
}

// newClusterRepository creates a new cluster redis operator via redis.ClusterOptions
func newClusterRepository(options *redis.ClusterOptions) cache.Operator {
	return &CacheRepository{cluster.Operator{ClusterClient: redis.NewClusterClient(options)}}
}

// InitCacheOperatorIfNecessary initializes the cache operator of repository
func InitCacheOperatorIfNecessary() error {
	if nil == config.GetConfigs() {
		log.Infosf("Cannnot found the config of service, skip init cache operator!")
		return nil
	}
	addressingConfig := &config.GetConfigs().Addressing

	// skip if addressing is disabled
	if !addressingConfig.Enable {
		log.Debugsf("addressing is disabled, skip init cache operator")
		return nil
	}

	if addressingConfig.SyncConfigWithServer {
		// request to GLS server to get cache config
		serviceConfig := config.GetConfigs().Service
		client := mesh.NewMeshClient()
		request := mesh.NewMeshRequest(nil)
		responseData := glsdef.RedisResponse{}
		request.WithOptions(
			mesh.WithTopicTypeBusiness(),                            // mark topic type to TRN
			mesh.WithORG(serviceConfig.Org),                         // org id
			mesh.WithWorkspace(serviceConfig.Wks),                   // workspace
			mesh.WithEnvironment(serviceConfig.Env),                 // environment
			mesh.WithSU(serviceConfig.Su),                           // su
			mesh.WithVersion(addressingConfig.TopicVersionOfServer), // destination event version
			mesh.WithEventID(addressingConfig.TopicIDOfServer),      // destination event id
			mesh.WithTimeout(30*time.Second),
		)

		// sync call
		_, err := client.SyncCall(context.Background(), request, &responseData)
		if nil != err {
			return err
		}

		if responseData.ErrorCode != 0 {
			return errors.Errorf(constant.SystemInternalError,
				"Failed to get cache config, error:%s", responseData.ErrorMsg)
		}

		addressingConfig.Cache.Type = responseData.Response.Type
		addressingConfig.Cache.Addr = responseData.Response.Addr
		addressingConfig.Cache.Password = responseData.Response.Password
		addressingConfig.Cache.Readonly = responseData.Response.Readonly
		addressingConfig.Cache.PoolNum = responseData.Response.Poolnum
	}
	var operator cache.Operator
	if len(strings.Split(addressingConfig.Cache.Addr, ",")) > 1 {
		// cluster
		operator = newClusterRepository(&redis.ClusterOptions{
			Addrs:    strings.Split(addressingConfig.Cache.Addr, ","),
			Password: addressingConfig.Cache.Password,
			PoolSize: addressingConfig.Cache.PoolNum,
			ReadOnly: addressingConfig.Cache.Readonly,
		})
	} else {
		// singleton
		operator = newSingletonRepository(&redis.Options{
			Addr:     addressingConfig.Cache.Addr,
			Password: addressingConfig.Cache.Password,
			PoolSize: addressingConfig.Cache.PoolNum,
		})
	}
	//if strings.EqualFold("cluster", addressingConfig.Cache.Type) {
	//	if len(strings.Split(addressingConfig.Cache.Addr, ",")) <= 1 {
	//		return errors.Errorf(constant.SystemInternalError,
	//			"Failed to create redis culster instance for lookup, the address[%s] is invalid!", addressingConfig.Cache.Addr)
	//	}
	//	operator = newClusterRepository(&redis.ClusterOptions{
	//		Addrs:    strings.Split(addressingConfig.Cache.Addr, ","),
	//		Password: addressingConfig.Cache.Password,
	//		PoolSize: addressingConfig.Cache.PoolNum,
	//		ReadOnly: addressingConfig.Cache.Readonly,
	//	})
	//} else {
	//	if len(strings.Split(addressingConfig.Cache.Addr, ",")) != 1 {
	//		return errors.Errorf(constant.SystemInternalError,
	//			"Failed to create redis singleton instance for lookup, the address[%s] is invalid!", addressingConfig.Cache.Addr)
	//	}
	//	operator = newSingletonRepository(&redis.Options{
	//		Addr:     addressingConfig.Cache.Addr,
	//		Password: addressingConfig.Cache.Password,
	//		PoolSize: addressingConfig.Cache.PoolNum,
	//	})
	//}

	cache.AddressingCacheOperator = operator
	log.Infosf("Successfully to init addressing cache operator, cache type is:%s", addressingConfig.Cache.Type)
	return nil
}

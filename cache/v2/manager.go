package cache

import (
	"git.multiverse.io/eventkit/kit/cache/v2/redis"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	redis2 "github.com/go-redis/redis/v8"
	"sync"
)

// ConnectionPool is used to define the operate of connection pool cache
type ConnectionPool interface {
	Get(aliasName string) (cp interface{}, ok bool)
	// Delete May close Cache before delete the connection pool by alias name
	Delete(aliasName string)
	InitCache(aliasName string, cacheConfig *config.Cache) *errors.Error
}

var _connectionPoolCache ConnectionPool

func rotateCacheConfigWhenConfigChanged(oldConfig *config.ServiceConfigs, newConfig *config.ServiceConfigs) error {
	if nil == newConfig {
		return nil
	}
	if err := Rotate(newConfig.Cache); nil != err {
		return err
	}
	return nil
}

// InitCache creates a connection pool cache
func InitCacheManagerForRedis(cacheConfigs map[string]config.Cache) *errors.Error {
	_connectionPoolCache = redis.NewRedisConnectionPoolCache()
	if err := Rotate(cacheConfigs); nil != err {
		return err
	}
	// register config change hook function into config manager for redis
	config.RegisterConfigOnChangeHookFunc("CacheManagerForRedis", rotateCacheConfigWhenConfigChanged, true)
	return nil
}

type _cachePools struct {
	sync.RWMutex
	Cache               map[string]suCacheConnectionPools
	CurrentCacheConfigs map[string]*config.Cache
}

type suCacheConnectionPools struct {
	DefaultPoolName string
	PoolNameMapping map[string]string
}

var cachePools = new(_cachePools)

func groupCacheConfig(cacheConfigs map[string]*config.Cache) map[string][]*config.Cache {
	for k, cache := range cacheConfigs {
		cache.Name = k
	}

	groupedCacheConfigs := make(map[string][]*config.Cache, 0)
	for _, cacheConfig := range cacheConfigs {
		if len(cacheConfig.Su) != 0 {
			if _, ok := groupedCacheConfigs[cacheConfig.Su]; !ok {
				groupedCacheConfigs[cacheConfig.Su] = make([]*config.Cache, 0)
			}
			groupedCacheConfigs[cacheConfig.Su] = append(groupedCacheConfigs[cacheConfig.Su], cacheConfig)
		}
	}

	for _, groupedCache := range groupedCacheConfigs {
		if len(groupedCache) == 1 {
			groupedCache[0].Default = true
		}
	}
	return groupedCacheConfigs
}

func CheckSuCacheConfigs(suCacheConfigs map[string][]*config.Cache) *errors.Error {
	for su, cacheConfigs := range suCacheConfigs {
		hasDefault := false
		if len(cacheConfigs) > 1 {
			for _, cacheConfig := range cacheConfigs {
				if cacheConfig.Default {
					if hasDefault {
						return errors.Errorf(constant.SystemInternalError, "only accept one default cache in the SU:[%s]", su)
					}
					hasDefault = true
				}
			}

			if !hasDefault {
				return errors.Errorf(constant.SystemInternalError, "Cannot found the default cache in the SU:[%s]", su)
			}
		}
	}

	return nil
}

func Rotate(cacheConfigs map[string]config.Cache) (re *errors.Error) {
	toRotateConfigs := make(map[string]*config.Cache)

	for su, cfg := range cacheConfigs {
		newCfg := cfg.Clone()
		toRotateConfigs[su] = &newCfg
	}
	return _rotate(toRotateConfigs)
}

func _rotate(cacheConfigs map[string]*config.Cache) (re *errors.Error) {
	groupedCacheConfigs := groupCacheConfig(cacheConfigs)

	if e := CheckSuCacheConfigs(groupedCacheConfigs); nil != e {
		return e
	}

	cachePools.Lock()
	defer cachePools.Unlock()
	defer func() {
		log.Infosf("The current Cache pool config:%++v", cachePools)
	}()
	defer func() {
		cachePools.CurrentCacheConfigs = cacheConfigs
	}()

	if 0 == len(cachePools.CurrentCacheConfigs) {
		// 1. Initialize the Cache, each SU will have a default Cache pool
		// * If a SU in the Cache configuration has only one Cache, then it will become the default Cache pool
		// * If a SU is configured with multiple Caches in the Cache configuration, only one of them must be set as default
		cachePools.Cache = make(map[string]suCacheConnectionPools)
		for su, groupCache := range groupedCacheConfigs {
			cachePools.Cache[su] = generateSuCachePools(su, groupCache)
		}

		// init all caches
		for _, groupedCache := range groupedCacheConfigs {
			for _, cacheConfig := range groupedCache {
				if err := _connectionPoolCache.InitCache(cacheConfig.Name, cacheConfig); nil != err {
					log.Errorsf("Failed to InitCache[config=[%++v]], error:%++v", *cacheConfig, err)
					return err
				}
			}
		}
		return nil
	}

	// analysis：
	//  1.Need to create a new connection pool
	//  2.The connection pool that needs to be deleted
	//  3.Connection pool that needs to be reconnected
	needToAddCacheConfigs := []*config.Cache{}
	needToDeleteCacheConfigs := []*config.Cache{}
	needToUpdateCacheConfigs := []*config.Cache{}
	currentCacheConfigs := cachePools.CurrentCacheConfigs
	for su, cacheConfig := range currentCacheConfigs {
		if _, ok := cacheConfigs[su]; !ok {
			needToDeleteCacheConfigs = append(needToDeleteCacheConfigs, cacheConfig)
		}
	}

	for su, cacheConfig := range cacheConfigs {
		if _, ok := currentCacheConfigs[su]; !ok {
			needToAddCacheConfigs = append(needToAddCacheConfigs, cacheConfig)
		}
	}

	for su, cacheConfig := range currentCacheConfigs {
		if newConfig, ok := cacheConfigs[su]; ok {
			if !cacheConfig.EqualsWithoutTopics(newConfig) || !arrayEqualsIgnoreOrder(cacheConfig.Topics, newConfig.Topics) {
				needToUpdateCacheConfigs = append(needToUpdateCacheConfigs, newConfig)
			}
		}
	}

	// init new cache connection pool
	for _, cacheConfig := range needToAddCacheConfigs {
		if err := _connectionPoolCache.InitCache(cacheConfig.Name, cacheConfig); nil != err {
			log.Errorsf("Failed to InitCache[config=[%++v]], error:%++v", *cacheConfig, err)
			return err
		}
	}

	// delete unused cache connection pool
	for _, cacheConfig := range needToDeleteCacheConfigs {
		_connectionPoolCache.Delete(cacheConfig.Name)
	}

	// check all need reonnection
	for _, cacheConfig := range needToUpdateCacheConfigs {
		if !cacheConfig.EqualsWithoutTopics(currentCacheConfigs[cacheConfig.Name]) {
			_connectionPoolCache.Delete(cacheConfig.Name)
			if err := _connectionPoolCache.InitCache(cacheConfig.Name, cacheConfig); nil != err {
				log.Errorsf("Failed to InitCache[config=[%++v]], error:%++v", *cacheConfig, err)
				return err
			}
		}
	}

	// According to groupedCacheConfigs, update the existing configuration,
	// but do not re-create or destroy the cache connection
	for su, _ := range cachePools.Cache {
		if _, ok := groupedCacheConfigs[su]; !ok {
			// Delete the existing SU cache connection pool
			delete(cachePools.Cache, su)
		}
	}

	// insert or update the su cache connection pools information!
	for su, cacheConfigs := range groupedCacheConfigs {
		cachePools.Cache[su] = generateSuCachePools(su, cacheConfigs)
	}

	return nil
}

func generateSuCachePools(su string, cachesConfigs []*config.Cache) suCacheConnectionPools {
	suCachePools := suCacheConnectionPools{}
	defaultName := ""
	suCachePools.PoolNameMapping = make(map[string]string, 0)
	for _, cacheConfig := range cachesConfigs {
		if cacheConfig.Default {
			defaultName = cacheConfig.Name
		}
		for _, topic := range cacheConfig.Topics {
			suCachePools.PoolNameMapping[topic] = cacheConfig.Name
		}
	}
	suCachePools.DefaultPoolName = defaultName
	return suCachePools
}

func arrayEqualsIgnoreOrder(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for _, e := range a1 {
		isExistsInOtherArray := false
		for _, oe := range a2 {
			if e == oe {
				isExistsInOtherArray = true
				break
			}
		}
		if !isExistsInOtherArray {
			return false
		}
	}
	return true
}

func getCP(su, topicID string) (interface{}, *errors.Error) {
	cachePools.RLock()
	defer cachePools.RUnlock()

	// 1. First determine whether there is a corresponding su cached connection pool, if not, an error will be reported
	if 0 == len(cachePools.Cache) {
		return nil, errors.Errorf(constant.SystemInternalError,
			"The cache doesn't have initialized,please check!")
	}

	suCachePools, ok := cachePools.Cache[su]
	if !ok || nil == suCachePools.PoolNameMapping {
		return nil, errors.Errorf(constant.SystemInternalError,
			"Cannot found the Cache connection pool[su=%s]", su)
	}

	// 2. Under SU, match the corresponding Cache pool according to topicID,
	// if it fails to match, assign the current SU default Cache pool and return
	cachePoolName, ok := suCachePools.PoolNameMapping[topicID]

	if !ok && "" == suCachePools.DefaultPoolName {
		return nil, errors.Errorf(constant.SystemInternalError,
			"No suitable Cache connection pool[su=%s,topic id=%s]", su, topicID)
	}

	var aliasName string
	if ok {
		aliasName = cachePoolName

	} else {
		aliasName = suCachePools.DefaultPoolName
	}

	_, ok = cachePools.CurrentCacheConfigs[aliasName]
	if !ok {
		return nil, errors.Errorf(constant.SystemInternalError, "Cannot found the Cache config[aliasName=%s]", aliasName)
	}
	o, ok := _connectionPoolCache.Get(aliasName)
	if !ok {
		return nil, errors.Errorf(constant.SystemInternalError, "Faield to get Engine with key:%s", aliasName)
	}

	return o, nil
}

func GetRedisClient(su string, topicIDs ...string) (redis2.UniversalClient, *errors.Error) {
	topicID := ""
	if len(topicIDs) > 0 {
		topicID = topicIDs[0]
	}
	o, err := getCP(su, topicID)
	if nil != err {
		return nil, err
	}

	return o.(redis2.UniversalClient), nil
}

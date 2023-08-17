package db

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/db/beego"
	xormConnectionPoolCache "git.multiverse.io/eventkit/kit/db/xorm"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/xormplus/xorm"
	"sync"
)

// ConnectionPoolCache is used to define the operate of connection pool cache
type ConnectionPoolCache interface {
	Get(driverName, aliasName string) (cp interface{}, ok bool)
	// Delete May close DB before delete the connection pool by alias name
	Delete(driverName, aliasName string)
	InitDatabase(dbAliasName string, dbConfig *config.Db) *errors.Error
}

var _connectionPoolCache ConnectionPoolCache

func rotateDBConfigWhenConfigChanged(oldConfig *config.ServiceConfigs, newConfig *config.ServiceConfigs) error {
	if nil == newConfig {
		return nil
	}
	if err := Rotate(newConfig.Db); nil != err {
		return err
	}
	return nil
}

// InitDBManagerForXorm creates a connection pool cache for XORM
func InitDBManagerForXorm(dbConfigs map[string]config.Db) *errors.Error {
	_connectionPoolCache = xormConnectionPoolCache.NewXormConnectionPoolCache()
	if err := Rotate(dbConfigs); nil != err {
		return err
	}
	// register config change hook function into config manager for XORM
	config.RegisterConfigOnChangeHookFunc("DBManagerForXorm", rotateDBConfigWhenConfigChanged, true)
	return nil
}

// InitDBManagerForBeegoOrmer creates a connection pool cache for beego ormer
func InitDBManagerForBeegoOrmer(dbConfigs map[string]config.Db) *errors.Error {
	_connectionPoolCache = beego.NewBeegoConnectionPoolCache()
	if err := Rotate(dbConfigs); nil != err {
		return err
	}
	// register config change hook function into config manager for Beego Ormer
	config.RegisterConfigOnChangeHookFunc("DBManagerForBeegoOrmer", rotateDBConfigWhenConfigChanged, true)
	return nil
}

type _dbPoolsCache struct {
	sync.RWMutex
	Cache            map[string]suDBPools
	CurrentDBConfigs map[string]*config.Db
}

type suDBPools struct {
	DefaultPoolName string
	PoolNameMapping map[string]string
}

var dbPoolsCache = new(_dbPoolsCache)

func groupDBConfig(dbConfigs map[string]*config.Db) map[string][]*config.Db {
	for k, db := range dbConfigs {
		db.Name = k
	}

	groupedDbConfigs := make(map[string][]*config.Db, 0)
	for _, dbConfig := range dbConfigs {
		if len(dbConfig.Su) != 0 {
			if _, ok := groupedDbConfigs[dbConfig.Su]; !ok {
				groupedDbConfigs[dbConfig.Su] = make([]*config.Db, 0)
			}
			groupedDbConfigs[dbConfig.Su] = append(groupedDbConfigs[dbConfig.Su], dbConfig)
		}
	}

	for _, dbConfigs := range groupedDbConfigs {
		if len(dbConfigs) == 1 {
			dbConfigs[0].Default = true
		}
	}
	return groupedDbConfigs
}

func CheckSuDbConfigs(suDbConfigs map[string][]*config.Db) *errors.Error {
	for su, dbConfigs := range suDbConfigs {
		hasDefault := false
		if len(dbConfigs) > 1 {
			for _, dbConfig := range dbConfigs {
				if dbConfig.Default {
					if hasDefault {
						return errors.Errorf(constant.SystemInternalError, "only accept one default database in the SU:[%s]", su)
					}
					hasDefault = true
				}
			}

			if !hasDefault {
				return errors.Errorf(constant.SystemInternalError, "Cannot found the default database in the SU:[%s]", su)
			}
		}
	}

	return nil
}

func Rotate(dbConfigs map[string]config.Db) (re *errors.Error) {
	toRotateConfigs := make(map[string]*config.Db)

	for su, cfg := range dbConfigs {
		newCfg := cfg.Clone()
		toRotateConfigs[su] = &newCfg
	}
	return _rotate(toRotateConfigs)
}

func _rotate(dbConfigs map[string]*config.Db) (re *errors.Error) {
	groupedDbConfigs := groupDBConfig(dbConfigs)

	if e := CheckSuDbConfigs(groupedDbConfigs); nil != e {
		return e
	}

	dbPoolsCache.Lock()
	defer dbPoolsCache.Unlock()
	defer func() {
		log.Infosf("The current DB pool cache config:%++v", dbPoolsCache)
	}()
	defer func() {
		dbPoolsCache.CurrentDBConfigs = dbConfigs
	}()

	if 0 == len(dbPoolsCache.CurrentDBConfigs) {
		// 1. Initialize the DB, each SU will have a default DB pool
		// * If a SU in the DB configuration has only one DB, then it will become the default DB pool
		// * If a SU is configured with multiple DBs in the DB configuration, only one of them must be set as default
		dbPoolsCache.Cache = make(map[string]suDBPools)
		for su, dbConfigs := range groupedDbConfigs {
			dbPoolsCache.Cache[su] = generateSuDBPools(su, dbConfigs)
		}

		// init all databases
		for _, dbConfigs := range groupedDbConfigs {
			for _, dbConfig := range dbConfigs {
				if err := _connectionPoolCache.InitDatabase(dbConfig.Name, dbConfig); nil != err {
					log.Errorsf("Failed to InitDatabase[config=[%++v]], error:%++v", *dbConfig, err)
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
	needToAddDbConfigs := []*config.Db{}
	needToDeleteDbConfigs := []*config.Db{}
	needToUpdateDbConfigs := []*config.Db{}
	currentDbConfigs := dbPoolsCache.CurrentDBConfigs
	for su, dbConfig := range currentDbConfigs {
		if _, ok := dbConfigs[su]; !ok {
			needToDeleteDbConfigs = append(needToDeleteDbConfigs, dbConfig)
		}
	}

	for su, dbConfig := range dbConfigs {
		if _, ok := currentDbConfigs[su]; !ok {
			needToAddDbConfigs = append(needToAddDbConfigs, dbConfig)
		}
	}

	for su, dbConfig := range currentDbConfigs {
		if newConfig, ok := dbConfigs[su]; ok {
			//if !isDbConfigEqual(dbConfig, newConfig) || !arrayEqualsIgnoreOrder(dbConfig.Topics, newConfig.Topics) {
			if !dbConfig.EqualsWithoutTopics(newConfig) || !arrayEqualsIgnoreOrder(dbConfig.Topics, newConfig.Topics) {
				needToUpdateDbConfigs = append(needToUpdateDbConfigs, newConfig)
			}
		}
	}

	// init new database connection pool
	for _, dbConfig := range needToAddDbConfigs {
		if err := _connectionPoolCache.InitDatabase(dbConfig.Name, dbConfig); nil != err {
			log.Errorsf("Failed to InitDatabase[config=[%++v]], error:%++v", *dbConfig, err)
			return err
		}
	}

	// delete unused database connection pool
	for _, dbConfig := range needToDeleteDbConfigs {
		_connectionPoolCache.Delete(dbConfig.Type, dbConfig.Name)
	}

	// check all need reonnection
	for _, dbConfig := range needToUpdateDbConfigs {
		//if !isDbConfigEqual(currentDbConfigs[dbConfig.Name], dbConfig) {
		if !dbConfig.EqualsWithoutTopics(currentDbConfigs[dbConfig.Name]) {
			_connectionPoolCache.Delete(dbConfig.Type, dbConfig.Name)
			if err := _connectionPoolCache.InitDatabase(dbConfig.Name, dbConfig); nil != err {
				log.Errorsf("Failed to InitDatabase[config=[%++v]], error:%++v", *dbConfig, err)
				return err
			}
		}
	}

	// According to gourpedDbConfigs, update the existing configuration,
	// but do not re-create or destroy the database connection
	for su, _ := range dbPoolsCache.Cache {
		if _, ok := groupedDbConfigs[su]; !ok {
			// Delete the existing SU database connection pool
			delete(dbPoolsCache.Cache, su)
		}
	}

	// insert or update the su db connection pools information!
	for su, dbConfigs := range groupedDbConfigs {
		dbPoolsCache.Cache[su] = generateSuDBPools(su, dbConfigs)
	}

	return nil
}

func generateSuDBPools(su string, dbConfigs []*config.Db) suDBPools {
	suDBPools := suDBPools{}
	defaultName := ""
	suDBPools.PoolNameMapping = make(map[string]string, 0)
	for _, dbConfig := range dbConfigs {
		if dbConfig.Default {
			defaultName = dbConfig.Name
		}
		for _, topic := range dbConfig.Topics {
			suDBPools.PoolNameMapping[topic] = dbConfig.Name
		}
	}
	suDBPools.DefaultPoolName = defaultName
	return suDBPools
}

//
//func isDbConfigEqual(f *config.Db, s *config.Db) bool {
//	return f.Type == s.Type &&
//		f.Su == s.Su &&
//		f.Default == s.Default &&
//		f.Addr == s.Addr &&
//		f.User == s.User &&
//		f.Password == s.Password &&
//		f.Database == s.Database &&
//		f.Params == s.Params &&
//		f.Debug == s.Debug &&
//		f.Pool.MaxIdleConns == s.Pool.MaxIdleConns &&
//		f.Pool.MaxOpenConns == s.Pool.MaxOpenConns &&
//		f.Pool.DefaultQueryLimit == s.Pool.DefaultQueryLimit &&
//		f.Pool.MaxLimitValue == s.Pool.MaxLimitValue &&
//		f.Pool.MaxLifeValue == s.Pool.MaxLifeValue
//
//}

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
	dbPoolsCache.RLock()
	defer dbPoolsCache.RUnlock()

	// 1. First determine whether there is a corresponding su database connection pool, if not, an error will be reported
	if 0 == len(dbPoolsCache.Cache) {
		return nil, errors.Errorf(constant.SystemInternalError,
			"The database source doesn't have initialized,please check!")
	}

	suDBPools, ok := dbPoolsCache.Cache[su]
	if !ok || nil == suDBPools.PoolNameMapping {
		return nil, errors.Errorf(constant.SystemInternalError,
			"Cannot found the DB connection pool[su=%s]", su)
	}

	// 2. Under SU, match the corresponding DB pool according to topicID,
	// if it fails to match, assign the current SU default DB pool and return
	dbPoolName, ok := suDBPools.PoolNameMapping[topicID]

	if !ok && "" == suDBPools.DefaultPoolName {
		return nil, errors.Errorf(constant.SystemInternalError,
			"No suitable DB connection pool[su=%s,topic id=%s]", su, topicID)
	}

	var aliasName string
	if ok {
		aliasName = dbPoolName

	} else {
		aliasName = suDBPools.DefaultPoolName
	}

	cfg, ok := dbPoolsCache.CurrentDBConfigs[aliasName]
	if !ok {
		return nil, errors.Errorf(constant.SystemInternalError, "Cannot found the DB config[aliasName=%s]", aliasName)
	}
	o, ok := _connectionPoolCache.Get(cfg.Type, aliasName)
	if !ok {
		return nil, errors.Errorf(constant.SystemInternalError, "Faield to get Engine with key:%s", aliasName)
	}

	return o, nil
}

func GetXormEngine(su string, topicIDs ...string) (*xorm.Engine, *errors.Error) {
	topicID := ""
	if len(topicIDs) > 0 {
		topicID = topicIDs[0]
	}
	o, err := getCP(su, topicID)
	if nil != err {
		return nil, err
	}

	return o.(*xorm.Engine), nil
}

func GetBeegoOrmer(su string, topicIDs ...string) (orm.Ormer, *errors.Error) {
	topicID := ""
	if len(topicIDs) > 0 {
		topicID = topicIDs[0]
	}
	o, err := getCP(su, topicID)
	if nil != err {
		return nil, err
	}
	return o.(orm.Ormer), nil
}

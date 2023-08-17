package xorm

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/db/util"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/xormplus/xorm"
)

// NewXormConnectionPoolCache creates a XORM connection pool cache
func NewXormConnectionPoolCache() *_cache {
	var enginesCache = &_cache{
		cache: make(map[string]*xorm.Engine),
	}

	return enginesCache
}

type _cache struct {
	sync.RWMutex
	cache map[string]*xorm.Engine
}

func (e *_cache) Get(driverName, aliasName string) (eg interface{}, ok bool) {
	e.RLock()
	defer e.RUnlock()
	if len(e.cache) > 0 {
		eg, ok = e.cache[aliasName]
	}

	return
}

func (e *_cache) set(driverName, aliasName string, eg *xorm.Engine) {
	e.Lock()
	defer e.Unlock()
	if nil == e.cache {
		e.cache = make(map[string]*xorm.Engine)
	}

	e.cache[aliasName] = eg
}

func (e *_cache) Delete(driverName, aliasName string) {
	eg, ok := e.Get(driverName, aliasName)
	if ok {
		// close DB
		eg.(*xorm.Engine).Close()
	}
	if nil != e.cache {
		delete(e.cache, aliasName)
	}
}

func (e *_cache) InitDatabase(aliasName string, dbConfig *config.Db) (err *errors.Error) {
	addr := dbConfig.Addr
	userName := dbConfig.User
	password := dbConfig.Password
	dbName := dbConfig.Database
	log.Infosf("Start init database[%++v] for Xorm", *dbConfig)
	var dsn string
	var gerr *errors.Error
	if strings.EqualFold("mysql", dbConfig.Type) {
		dsn = userName + ":" + password + "@tcp(" + addr + ")/" + dbName
		dsn, gerr = util.AppendSSLConnectionStringParamForMysqlIfNecessary(aliasName, dsn, dbConfig.Params, dbConfig.ServerCACertFile, dbConfig.ClientCertFile, dbConfig.ClientPriKeyFile)
		if nil != gerr {
			return gerr
		}
	} else if strings.EqualFold("postgres", dbConfig.Type) {
		addrs := strings.Split(addr, ":")
		if len(addrs) < 2 {
			return errors.Errorf(constant.SystemInternalError, "Invalid address of database:%s", addr)
		}
		host := addrs[0]
		port := addrs[1]
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, userName, password, dbName)
		dsn, gerr = util.AppendSSLConnectionStringParamForPostgresIfNecessary(dsn, dbConfig.Params, dbConfig.ServerCACertFile, dbConfig.ClientCertFile, dbConfig.ClientPriKeyFile)
		if nil != gerr {
			return gerr
		}
	} else {
		return errors.Errorf(constant.SystemInternalError, "unsupported db type: %s", dbConfig.Type)
	}

	engine, xerr := xorm.NewEngine(dbConfig.Type, dsn)

	if xerr != nil {
		return errors.Wrap(constant.SystemInternalError, xerr, 0)
	}

	if xerr := engine.Ping(); nil != xerr {
		return errors.Wrap(constant.SystemInternalError, xerr, 0)
	}

	maxIdleConns := dbConfig.Pool.MaxIdleConns
	maxOpenConns := dbConfig.Pool.MaxOpenConns
	maxIdleTime := dbConfig.Pool.MaxIdleTime
	maxLifeValue := dbConfig.Pool.MaxLifeValue
	showSql := dbConfig.Debug

	engine.SetMaxIdleConns(maxIdleConns)
	engine.SetMaxOpenConns(maxOpenConns)
	engine.DB().SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second)
	engine.SetConnMaxLifetime(time.Duration(maxLifeValue) * time.Second)
	engine.ShowSQL(showSql)

	if dbConfig.DBTimeZone != "" {
		tz, tzerr := time.LoadLocation(dbConfig.DBTimeZone)
		if nil != tzerr {
			return errors.Errorf(constant.SystemInternalError, "Failed to load timezone: %++v", tzerr)
		}

		// set timezone of database
		engine.SetTZDatabase(tz)
	}
	e.set(dbConfig.Type, aliasName, engine)
	log.Infosf("Successfully to init database[%++v] for Xorm", *dbConfig)
	return nil
}

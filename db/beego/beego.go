package beego

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/db/util"
	"git.multiverse.io/eventkit/kit/handler/config"
	"github.com/beego/beego/v2/adapter/orm"
)

// NewBeegoConnectionPoolCache creates a Beego connection pool cache
func NewBeegoConnectionPoolCache() *cache {
	var enginesCache = &cache{
		cache: make(map[string]*sql.DB),
	}

	return enginesCache
}

type cache struct {
	sync.RWMutex
	cache map[string]*sql.DB
}

func (e *cache) Get(driverName, aliasName string) (ormer interface{}, ok bool) {
	e.RLock()
	defer e.RUnlock()
	if len(e.cache) > 0 {
		_, iOk := e.cache[aliasName]
		if iOk {
			o := orm.NewOrm()
			o.Using(aliasName)

			return o, true
		}

		return nil, false
	}

	return nil, false
}

func (e *cache) set(driverName, aliasName string, db *sql.DB) {
	e.Lock()
	defer e.Unlock()
	if nil == e.cache {
		e.cache = make(map[string]*sql.DB)
	}

	e.cache[aliasName] = db
}

func (e *cache) Delete(driverName, aliasName string) {
	if db, ok := e.cache[aliasName]; ok {
		db.Close()
	}
	if nil != e.cache {
		delete(e.cache, aliasName)
	}
}

func (e *cache) InitDatabase(aliasName string, dbConfig *config.Db) (err *errors.Error) {
	addr := dbConfig.Addr
	userName := dbConfig.User
	password := dbConfig.Password
	dbName := dbConfig.Database
	maxIdleConns := dbConfig.Pool.MaxIdleConns
	maxOpenConns := dbConfig.Pool.MaxOpenConns
	maxIdleTime := dbConfig.Pool.MaxIdleTime
	maxLifeValue := dbConfig.Pool.MaxLifeValue

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

	db, oerr := sql.Open(dbConfig.Type, dsn)
	if nil != oerr {
		return errors.Errorf(constant.SystemInternalError, "Failed to open database set: %++v", err)
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(maxLifeValue) * time.Second)

	if oeer := orm.AddAliasWthDB(aliasName, dbConfig.Type, db); nil != oeer {
		return errors.Errorf(constant.SystemInternalError, "Failed to add alias with DB: %++v", oeer)
	}
	if dbConfig.DBTimeZone != "" {
		tz, tzerr := time.LoadLocation(dbConfig.DBTimeZone)
		if nil != tzerr {
			return errors.Errorf(constant.SystemInternalError, "Failed to load timezone: %++v", tzerr)
		}
		orm.SetDataBaseTZ(aliasName, tz)
	}
	e.set(dbConfig.Type, aliasName, db)

	return nil
}

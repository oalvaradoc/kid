package datasource

import (
	"database/sql/driver"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type ATDriver struct {
	*Driver
}

func (a *ATDriver) OpenConnector(dsn string) (driver.Connector, error) {
	connector, err := a.Driver.OpenConnector(dsn)

	if nil != err {
		return nil, err
	}

	conn, _ := connector.(*Connector)
	switch strings.ToLower(a.DBType) {
	case "mysql":
		{
			cfg, err := mysql.ParseDSN(dsn)
			if nil != err {
				return nil, err
			}
			conn.Config = cfg
		}
	default:
		return nil, errors.New("unknown db type: " + a.DBType)
	}

	return conn, nil
}

package datasource

import (
	"context"
	"database/sql/driver"
	"git.multiverse.io/eventkit/kit/db/datasource/types"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"strings"
)

type Driver struct {
	DBType string
	Mode  types.TransactionMode
	OriginalDriver driver.Driver
}


func (d *Driver) Open(name string) (driver.Conn, error) {
	return d.OriginalDriver.Open(name)
}

func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	var originalConnector driver.Connector
	var err error
	originalConnector = &dsnConnector{
		dsn:    dsn,
		driver: d.OriginalDriver,
	}

	if driverCtx, ok := d.OriginalDriver.(driver.DriverContext); ok {
		originalConnector, err = driverCtx.OpenConnector(dsn)
		if nil != err {
			return nil, err
		}
	}

	switch strings.ToLower(d.DBType) {
		case "mysql": {
			cfg, err := mysql.ParseDSN(dsn)
			if nil != err {
				return nil, err
			}
			ret := &Connector{
				OriginalDriver: d.OriginalDriver,
				OriginalConnector: originalConnector,
				Config: cfg,
			}
			return ret, nil
		}
	default:
		return nil, errors.New("unknown db type: " + d.DBType)
	}
}


type dsnConnector struct{
	dsn string
	driver driver.Driver
}

func (d *dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return d.driver.Open(d.dsn)
}

func (d *dsnConnector) Driver() driver.Driver {
	return d.driver
}
package connector

import (
	"context"
	"database/sql/driver"

	"git.multiverse.io/eventkit/kit/db/datasource"
)

type ATConnector struct {
	*datasource.Connector
}

func (a *ATConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := a.Connector.Connect(ctx)

	if nil != err {
		return nil, err
	}

	return conn, nil
}

func (a *ATConnector) Driver() driver.Driver {
	return &datasource.ATDriver{
		Driver: a.OriginalConnector.Driver().(*datasource.Driver),
	}
}

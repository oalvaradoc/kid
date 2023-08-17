package datasource

import (
	"context"
	"database/sql/driver"
	"git.multiverse.io/eventkit/kit/db/datasource/types"
	"github.com/go-sql-driver/mysql"
	"sync"
)

type Connector struct {
	sync.Once
	OriginalDriver driver.Driver
	OriginalConnector driver.Connector
	Config *mysql.Config
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.OriginalConnector.Connect(ctx)

	if nil != err {
		return nil, err
	}

	return &Connection{
		TxnCtx:       &types.TransactionContexts{
			RoundImages: &types.RoundRecordImage{},
		},
		OriginalConn: conn,
		AutoCommit:   true,
	}, nil
}

func (c *Connector) Driver() driver.Driver {
	c.Do(func() {
		c.OriginalDriver = c.OriginalConnector.Driver()
	})

	return c.OriginalDriver
}
package datasource

import (
	"context"
	"database/sql/driver"
	"git.multiverse.io/eventkit/kit/db/datasource/types"
)

type Connection struct {
	TxnCtx       *types.TransactionContexts
	OriginalConn driver.Conn
	AutoCommit   bool
}


func (c *Connection) ResetSession(ctx context.Context) error {
	if conn, ok := c.OriginalConn.(driver.SessionResetter); ok {
		c.AutoCommit = true
		return conn.ResetSession(ctx)
	}

	return driver.ErrSkip
}

func (c *Connection) Prepare(query string) (driver.Stmt, error) {
	if s, err := c.OriginalConn.Prepare(query); nil != err {
		return nil, err
	} else {
		return s, nil
	}
}

func (c *Connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	conn, ok := c.OriginalConn.(driver.ConnPrepareContext)
	if !ok {
		return c.OriginalConn.Prepare(query)

	}
	if s, err := conn.PrepareContext(ctx, query); nil != err {
		return nil, err
	} else {
		return s, nil
	}

}

func (c *Connection) Exec(query string, args [] driver.Value) (driver.Result, error) {
	conn, ok := c.OriginalConn.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	ret, err := conn.Exec(query, args)
	if nil != err {
		return nil, err
	}
	return ret, nil

}

func (c *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {

	conn, ok := c.OriginalConn.(driver.ExecerContext)
	if !ok {
		values := make([]driver.Value, 0, len(args))

		for i := range args {
			values = append(values, args[i].Value)
		}

		return c.Exec(query, values)
	}

	ret, err := conn.ExecContext(ctx, query, args)
	if nil != err {
		return nil, err
	}
	return ret, nil

}

func (c *Connection) Query(query string, args []driver.Value) (driver.Rows, error) {
	conn, ok := c.OriginalConn.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	return conn.Query(query, args)
}

func (c *Connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	conn, ok := c.OriginalConn.(driver.QueryerContext)
	if !ok {
		values := make([]driver.Value, 0, len(args))

		for i := range args {
			values = append(values, args[i].Value)
		}
		return c.Query(query, values)
	}

	return conn.QueryContext(ctx, query, args)
}


// Begin creates a new transaction
func (c *Connection) Begin() (driver.Tx, error) {
	if tx, err := c.OriginalConn.Begin(); nil != err {
		return nil, err
	} else {
		return tx, nil
	}
}


func (c *Connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	conn, ok := c.OriginalConn.(driver.ConnBeginTx)

	if ok {
		return conn.BeginTx(ctx, opts)
	}

	return c.Begin()
}

func (c *Connection) Close() error {
	return c.OriginalConn.Close()
}
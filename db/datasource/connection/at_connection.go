package connection

import (
	"context"
	"database/sql/driver"

	"git.multiverse.io/eventkit/kit/db/datasource"
	"git.multiverse.io/eventkit/kit/db/datasource/executor"
	"git.multiverse.io/eventkit/kit/db/datasource/types"
)

type ATConnection struct {
	*datasource.Connection
}

func (a *ATConnection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	conn, ok := a.OriginalConn.(driver.ConnPrepareContext)
	if !ok {
		if s, err := a.OriginalConn.Prepare(query); err != nil {
			return nil, err
		} else {
			return &datasource.Stmt{Ctx: ctx, Sqlquery: query, Stmt: s}, nil
		}
	}
	if s, err := conn.PrepareContext(ctx, query); nil != err {
		return nil, err
	} else {
		return &datasource.Stmt{Ctx: ctx, Sqlquery: query, Stmt: s}, nil
	}
}

func (a *ATConnection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	conn, ok := a.OriginalConn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	ret, err := executor.ATExecuteWithNamedValue(ctx, query, args,
		func(ctx context.Context, query string, args []driver.NamedValue) (*types.ExecuteResult, error) {
			rows, err := conn.QueryContext(ctx, query, args)
			if err != nil {
				return nil, err
			}
			return &types.ExecuteResult{Rows: &rows}, nil
		})

	if nil != err {
		return nil, err
	}
	return *ret.Rows, nil
}

func (a *ATConnection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	conn, ok := a.OriginalConn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	ret, err := executor.ATExecuteWithNamedValue(ctx, query, args,
		func(ctx context.Context, query string, args []driver.NamedValue) (*types.ExecuteResult, error) {
			result, err := conn.ExecContext(ctx, query, args)
			if err != nil {
				return nil, err
			}
			return &types.ExecuteResult{Result: &result}, nil
		})

	if nil != err {
		return nil, err
	}
	return *ret.Result, nil
}

func (a *ATConnection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return a.Connection.BeginTx(ctx, opts)
}

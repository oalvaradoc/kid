package datasource

import (
	"context"
	"database/sql/driver"
	"errors"

	"git.multiverse.io/eventkit/kit/db/datasource/executor"
	"git.multiverse.io/eventkit/kit/db/datasource/types"
)

type Stmt struct {
	// conn *Conn
	// res   *DBResource
	// txCtx *types.TransactionContext
	Ctx      context.Context
	Sqlquery string
	Stmt     driver.Stmt
}

// Close closes the statement.
//
// As of Go 1.1, a Stmt will not be closed if it's in use
// by any queries.
//
// Drivers must ensure all network calls made by Close
// do not block indefinitely (e.g. apply a timeout).
func (s *Stmt) Close() error {
	// s.txCtx = nil
	return s.Stmt.Close()
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (s *Stmt) NumInput() int {
	return s.Stmt.NumInput()
}

// Query executes a query that may return rows, such as a
// SELECT.
//
// Deprecated: Drivers should implement StmtQueryContext instead (or additionally).
func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("should use QueryContext")
}

// QueryContext StmtQueryContext enhances the Stmt interface by providing Query with context.
// QueryContext executes a query that may return rows, such as a  SELECT.
// QueryContext must honor the context timeout and return when it is canceled.
func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	stmt, ok := s.Stmt.(driver.StmtQueryContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	ret, err := executor.ATExecuteWithNamedValue(context.Background(), s.Sqlquery, args,
		func(ctx context.Context, _ string, args []driver.NamedValue) (*types.ExecuteResult, error) {
			rows, err := stmt.QueryContext(ctx, args)
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

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// Deprecated: Drivers should implement StmtExecContext instead (or additionally).
func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("should use ExecContext")
}

// ExecContext executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	stmt, ok := s.Stmt.(driver.StmtExecContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	ret, err := executor.ATExecuteWithNamedValue(context.Background(), s.Sqlquery, args,
		func(ctx context.Context, _ string, args []driver.NamedValue) (*types.ExecuteResult, error) {
			result, err := stmt.ExecContext(ctx, args)
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

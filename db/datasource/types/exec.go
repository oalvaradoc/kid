package types

import (
	"context"
	"database/sql/driver"
)

type ExecuteResult struct {
	Rows   *driver.Rows
	Result *driver.Result
}

type CallBack func(ctx context.Context, query string, args []driver.NamedValue) (*ExecuteResult, error)

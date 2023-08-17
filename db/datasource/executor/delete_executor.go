package executor

import (
	"context"
	"database/sql/driver"

	"git.multiverse.io/eventkit/kit/db/datasource/types"
)

type DeleteExecutor struct {
	ATExecutor
}

func (e *DeleteExecutor) ExecContext(ctx context.Context, args []driver.NamedValue, f types.CallBack) (*types.ExecuteResult, error) {
	return f(ctx, e.query, args)
}

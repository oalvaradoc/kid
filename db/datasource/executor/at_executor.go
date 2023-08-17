package executor

import (
	"context"
	"database/sql/driver"
	"hash/crc32"

	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/db/datasource/types"
	"git.multiverse.io/eventkit/kit/db/util"
	"github.com/arana-db/parser"
	"github.com/arana-db/parser/ast"
)

var executorCache = make(map[uint32]Executor, 1)

type Executor interface {
	ExecContext(ctx context.Context, args []driver.NamedValue, f types.CallBack) (*types.ExecuteResult, error)
}

type ATExecutor struct {
	query string
	istmt *ast.InsertStmt
	ustmt *ast.UpdateStmt
	sstmt *ast.SelectStmt
	dstmt *ast.DeleteStmt
}

func (e *ATExecutor) ExecContext(ctx context.Context, args []driver.NamedValue, f types.CallBack) (*types.ExecuteResult, error) {
	return f(ctx, e.query, args)
}

func ATExecuteWithNamedValue(ctx context.Context, query string, args []driver.NamedValue, f types.CallBack) (*types.ExecuteResult, error) {
	crc := crc32.ChecksumIEEE([]byte(query))
	atexecutor, has := executorCache[crc]
	if !has {
		p := parser.New()
		stmtNodes, _, err := p.Parse(query, "", "")
		if err != nil {
			return nil, errors.Errorf("xxx", "not implement error")
		}

		if len(stmtNodes) != 1 {
			return nil, errors.Errorf("xxx", "not inplement multi")
		}

		switch stmt := stmtNodes[0].(type) {
		case *ast.InsertStmt:
			// parserCtx.SQLType = types.SQLTypeInsert
			// parserCtx.InsertStmt = stmt
			// parserCtx.ExecutorType = types.InsertExecutor

			// if stmt.IsReplace {
			// 	parserCtx.ExecutorType = types.ReplaceIntoExecutor
			// }
			// if len(stmt.OnDuplicate) != 0 {
			// 	parserCtx.SQLType = types.SQLTypeInsertOnDuplicateUpdate
			// 	parserCtx.ExecutorType = types.InsertOnDuplicateExecutor
			// }
			atexecutor = &InsertExecutor{ATExecutor{query: query, istmt: stmt}}
		case *ast.UpdateStmt:
			// parserCtx.SQLType = types.SQLTypeUpdate
			// parserCtx.UpdateStmt = stmt
			// parserCtx.ExecutorType = types.UpdateExecutor
			atexecutor = &UpdateExecutor{ATExecutor{query: query, ustmt: stmt}}
		case *ast.SelectStmt:
			// if stmt.LockInfo != nil && stmt.LockInfo.LockType == ast.SelectLockForUpdate {
			// 	parserCtx.SQLType = types.SQLTypeSelectForUpdate
			// 	parserCtx.SelectStmt = stmt
			// 	parserCtx.ExecutorType = types.SelectForUpdateExecutor
			// } else {
			// 	parserCtx.SQLType = types.SQLTypeSelect
			// 	parserCtx.SelectStmt = stmt
			// 	parserCtx.ExecutorType = types.SelectExecutor
			// }
			atexecutor = &SelectExecutor{ATExecutor{query: query, sstmt: stmt}}
		case *ast.DeleteStmt:
			// parserCtx.SQLType = types.SQLTypeDelete
			// parserCtx.DeleteStmt = stmt
			// parserCtx.ExecutorType = types.DeleteExecutor
			atexecutor = &DeleteExecutor{ATExecutor{query: query, dstmt: stmt}}

		}
		executorCache[crc] = atexecutor
	}

	return atexecutor.ExecContext(ctx, args, f)
}

func ATExecuteWithValue(ctx context.Context, query string, args []driver.Value, f types.CallBack) (*types.ExecuteResult, error) {
	return ATExecuteWithNamedValue(ctx, query, util.ValueToNamedValue(args), f)
}

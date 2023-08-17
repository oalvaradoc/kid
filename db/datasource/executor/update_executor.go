package executor

import (
	"bytes"
	"context"
	"database/sql/driver"

	"git.multiverse.io/eventkit/kit/db/datasource/types"
	"github.com/arana-db/parser/ast"
	"github.com/arana-db/parser/format"
	"github.com/arana-db/parser/model"
)

type UpdateExecutor struct {
	ATExecutor
}

func (e *UpdateExecutor) ExecContext(ctx context.Context, args []driver.NamedValue, f types.CallBack) (*types.ExecuteResult, error) {
	// selectSql, err := e.BuildSelectSql()
	// if err != nil {
	// 	return nil, nil, err
	// }
	// e.BuildBeforeImage()
	// e.con.ExecContext(ctx, selectSql, args)
	return f(ctx, e.query, args)
}

func (e *UpdateExecutor) BuildSelectSql() (string, error) {

	fields := make([]*ast.SelectField, 0, len(e.ustmt.List))

	// if undo.UndoConfig.OnlyCareUpdateColumns {
	// 	for _, column := range updateStmt.List {
	// 		fields = append(fields, &ast.SelectField{
	// 			Expr: &ast.ColumnNameExpr{
	// 				Name: column.Column,
	// 			},
	// 		})
	// 	}

	// 	// select indexes columns
	// 	tableName, _ := e.parserCtx.GetTableName()
	// 	metaData, err := datasource.GetTableCache(types.DBTypeMySQL).GetTableMeta(ctx, e.execContext.DBName, tableName)
	// 	if err != nil {
	// 		return "", nil, err
	// 	}
	// 	for _, columnName := range metaData.GetPrimaryKeyOnlyName() {
	// 		fields = append(fields, &ast.SelectField{
	// 			Expr: &ast.ColumnNameExpr{
	// 				Name: &ast.ColumnName{
	// 					Name: model.CIStr{
	// 						O: columnName,
	// 						L: columnName,
	// 					},
	// 				},
	// 			},
	// 		})
	// 	}
	// } else {
	fields = append(fields, &ast.SelectField{
		Expr: &ast.ColumnNameExpr{
			Name: &ast.ColumnName{
				Name: model.CIStr{
					O: "*",
					L: "*",
				},
			},
		},
	})
	// }

	selStmt := ast.SelectStmt{
		SelectStmtOpts: &ast.SelectStmtOpts{},
		From:           e.ustmt.TableRefs,
		Where:          e.ustmt.Where,
		Fields:         &ast.FieldList{Fields: fields},
		OrderBy:        e.ustmt.Order,
		Limit:          e.ustmt.Limit,
		TableHints:     e.ustmt.TableHints,
		LockInfo: &ast.SelectLockInfo{
			LockType: ast.SelectLockForUpdate,
		},
	}

	b := bytes.NewBuffer([]byte{})
	_ = selStmt.Restore(format.NewRestoreCtx(format.RestoreKeyWordUppercase, b))
	return b.String(), nil
}

func (e *UpdateExecutor) BuildBeforeImage() {

}

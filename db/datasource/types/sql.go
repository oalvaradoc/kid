package types

//go:generate stringer -type=SQLType
type SQLType int32

const (
	SQLTypeSelect = iota
	SQLTypeInsert
	SQLTypeUpdate
	SQLTypeDelete
	SQLTypeSelectForUpdate
	SQLTypeReplace
	SQLTypeTruncate
	SQLTypeCreate
	SQLTypeDrop
	SQLTypeLoad
	SQLTypeMerge
	SQLTypeShow
	SQLTypeAlter
	SQLTypeRename
	SQLTypeDump
	SQLTypeDebug
	SQLTypeExplain
	SQLTypeProcedure
	SQLTypeDesc
	SQLLastInsertID
	SQLSelectWithoutTable
	SQLCreateSequence
	SQLShowSequence
	SQLGetSequence
	SQLAlterSequence
	SQLDropSequence
	SQLTddlShow
	SQLTypeSet
	SQLTypeReload
	SQLTypeSelectUnion
	SQLTypeCreateTable
	SQLTypeDropTable
	SQLTypeAlterTable
	SQLTypeSavePoint
	SQLTypeSelectFromUpdate
	SQLTypeMultiDelete
	SQLTypeMultiUpdate
	SQLTypeCreateIndex
	SQLTypeDropIndex
	SQLTypeKill
	SQLTypeLockTables
	SQLTypeUnLockTables
	SQLTypeCheckTable
	SQLTypeSelectFoundRows
	SQLTypeInsertIgnore = iota + 57
	SQLTypeInsertOnDuplicateUpdate
	SQLTypeMulti = iota + 999
	SQLTypeUnknown
)

func (s SQLType) MarshalText() (text []byte, err error) {
	switch s {
	case SQLTypeSelect:
		return []byte("SELECT"), nil
	case SQLTypeInsert:
		return []byte("INSERT"), nil
	case SQLTypeUpdate:
		return []byte("UPDATE"), nil
	case SQLTypeDelete:
		return []byte("DELETE"), nil
	case SQLTypeSelectForUpdate:
		return []byte("SELECT_FOR_UPDATE"), nil
	case SQLTypeInsertOnDuplicateUpdate:
		return []byte("INSERT_ON_UPDATE"), nil
	case SQLTypeReplace:
		return []byte("REPLACE"), nil
	case SQLTypeTruncate:
		return []byte("TRUNCATE"), nil
	case SQLTypeCreate:
		return []byte("CREATE"), nil
	case SQLTypeDrop:
		return []byte("DROP"), nil
	case SQLTypeLoad:
		return []byte("LOAD"), nil
	case SQLTypeMerge:
		return []byte("MERGE"), nil
	case SQLTypeShow:
		return []byte("SHOW"), nil
	case SQLTypeAlter:
		return []byte("ALTER"), nil
	case SQLTypeRename:
		return []byte("RENAME"), nil
	case SQLTypeDump:
		return []byte("DUMP"), nil
	case SQLTypeDebug:
		return []byte("DEBUG"), nil
	case SQLTypeExplain:
		return []byte("EXPLAIN"), nil
	case SQLTypeDesc:
		return []byte("DESC"), nil
	case SQLTypeSet:
		return []byte("SET"), nil
	case SQLTypeReload:
		return []byte("RELOAD"), nil
	case SQLTypeSelectUnion:
		return []byte("SELECT_UNION"), nil
	case SQLTypeCreateTable:
		return []byte("CREATE_TABLE"), nil
	case SQLTypeDropTable:
		return []byte("DROP_TABLE"), nil
	case SQLTypeAlterTable:
		return []byte("ALTER_TABLE"), nil
	case SQLTypeSelectFromUpdate:
		return []byte("SELECT_FROM_UPDATE"), nil
	case SQLTypeMultiDelete:
		return []byte("MULTI_DELETE"), nil
	case SQLTypeMultiUpdate:
		return []byte("MULTI_UPDATE"), nil
	case SQLTypeCreateIndex:
		return []byte("CREATE_INDEX"), nil
	case SQLTypeDropIndex:
		return []byte("DROP_INDEX"), nil
	case SQLTypeMulti:
		return []byte("MULTI"), nil
	}
	return []byte("INVALID_SQLTYPE"), nil
}

func (s *SQLType) UnmarshalText(b []byte) error {
	switch string(b) {
	case "SELECT":
		*s = SQLTypeSelect
	case "INSERT":
		*s = SQLTypeInsert
	case "UPDATE":
		*s = SQLTypeUpdate
	case "DELETE":
		*s = SQLTypeDelete
	case "SELECT_FOR_UPDATE":
		*s = SQLTypeSelectForUpdate
	case "INSERT_ON_UPDATE":
		*s = SQLTypeInsertOnDuplicateUpdate
	case "REPLACE":
		*s = SQLTypeReplace
	case "TRUNCATE":
		*s = SQLTypeTruncate
	case "CREATE":
		*s = SQLTypeCreate
	case "DROP":
		*s = SQLTypeDrop
	case "LOAD":
		*s = SQLTypeLoad
	case "MERGE":
		*s = SQLTypeMerge
	case "SHOW":
		*s = SQLTypeShow
	case "ALTER":
		*s = SQLTypeAlter
	case "RENAME":
		*s = SQLTypeRename
	case "DUMP":
		*s = SQLTypeDump
	case "DEBUG":
		*s = SQLTypeDebug
	case "EXPLAIN":
		*s = SQLTypeExplain
	case "DESC":
		*s = SQLTypeDesc
	case "SET":
		*s = SQLTypeSet
	case "RELOAD":
		*s = SQLTypeReload
	case "SELECT_UNION":
		*s = SQLTypeSelectUnion
	case "CREATE_TABLE":
		*s = SQLTypeCreateTable
	case "DROP_TABLE":
		*s = SQLTypeDropTable
	case "ALTER_TABLE":
		*s = SQLTypeAlterTable
	case "SELECT_FROM_UPDATE":
		*s = SQLTypeSelectFromUpdate
	case "MULTI_DELETE":
		*s = SQLTypeMultiDelete
	case "MULTI_UPDATE":
		*s = SQLTypeMultiUpdate
	case "CREATE_INDEX":
		*s = SQLTypeCreateIndex
	case "DROP_INDEX":
		*s = SQLTypeDropIndex
	case "MULTI":
		*s = SQLTypeMulti
	}
	return nil
}

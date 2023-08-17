package types

import (
	"database/sql"
	"reflect"
)

// https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnType
type FieldType byte

const (
	FieldTypeDecimal FieldType = iota
	FieldTypeTiny
	FieldTypeShort
	FieldTypeLong
	FieldTypeFloat
	FieldTypeDouble
	FieldTypeNULL
	FieldTypeTimestamp
	FieldTypeLongLong
	FieldTypeInt24
	FieldTypeDate
	FieldTypeTime
	FieldTypeDateTime
	FieldTypeYear
	FieldTypeNewDate
	FieldTypeVarChar
	FieldTypeBit
)

const (
	FieldTypeJSON FieldType = iota + 0xf5
	FieldTypeNewDecimal
	FieldTypeEnum
	FieldTypeSet
	FieldTypeTinyBLOB
	FieldTypeMediumBLOB
	FieldTypeLongBLOB
	FieldTypeBLOB
	FieldTypeVarString
	FieldTypeString
	FieldTypeGeometry
)

type nullTime = sql.NullTime

var (
	ScanTypeFloat32   = reflect.TypeOf(float32(0))
	ScanTypeFloat64   = reflect.TypeOf(float64(0))
	ScanTypeInt8      = reflect.TypeOf(int8(0))
	ScanTypeInt16     = reflect.TypeOf(int16(0))
	ScanTypeInt32     = reflect.TypeOf(int32(0))
	ScanTypeInt64     = reflect.TypeOf(int64(0))
	ScanTypeNullFloat = reflect.TypeOf(sql.NullFloat64{})
	ScanTypeNullInt   = reflect.TypeOf(sql.NullInt64{})
	ScanTypeNullTime  = reflect.TypeOf(nullTime{})
	ScanTypeUint8     = reflect.TypeOf(uint8(0))
	ScanTypeUint16    = reflect.TypeOf(uint16(0))
	ScanTypeUint32    = reflect.TypeOf(uint32(0))
	ScanTypeUint64    = reflect.TypeOf(uint64(0))
	ScanTypeRawBytes  = reflect.TypeOf(sql.RawBytes{})
	ScanTypeUnknown   = reflect.TypeOf(new(interface{}))
)

type JDBCType int16

const (
	JDBCTypeBit                   JDBCType = -7
	JDBCTypeTinyInt               JDBCType = -6
	JDBCTypeSmallInt              JDBCType = 5
	JDBCTypeInteger               JDBCType = 4
	JDBCTypeBigInt                JDBCType = -5
	JDBCTypeFloat                 JDBCType = 6
	JDBCTypeReal                  JDBCType = 7
	JDBCTypeDouble                JDBCType = 8
	JDBCTypeNumberic              JDBCType = 2
	JDBCTypeDecimal               JDBCType = 3
	JDBCTypeChar                  JDBCType = 1
	JDBCTypeVarchar               JDBCType = 12
	JDBCTypeLongVarchar           JDBCType = -1
	JDBCTypeDate                  JDBCType = 91
	JDBCTypeTime                  JDBCType = 92
	JDBCTypeTimestamp             JDBCType = 93
	JDBCTypeBinary                JDBCType = -2
	JDBCTypeVarBinary             JDBCType = -3
	JDBCTypeLongVarBinary         JDBCType = -4
	JDBCTypeNull                  JDBCType = 0
	JDBCTypeOther                 JDBCType = 1111
	JDBCTypeJavaObject            JDBCType = 2000
	JDBCTypeDistinct              JDBCType = 2001
	JDBCTypeStruct                JDBCType = 2002
	JDBCTypeArray                 JDBCType = 2003
	JDBCTypeBlob                  JDBCType = 2004
	JDBCTypeClob                  JDBCType = 2005
	JDBCTypeRef                   JDBCType = 2006
	JDBCTypeDateLink              JDBCType = 70
	JDBCTypeBoolean               JDBCType = 16
	JDBCTypeRowID                 JDBCType = -8
	JDBCTypeNchar                 JDBCType = -15
	JDBCTypeNvarchar              JDBCType = -9
	JDBCTypeLongNvVarchar         JDBCType = -16
	JDBCTypeNclob                 JDBCType = 2011
	JDBCTypeSqlXML                JDBCType = 2009
	JDBCTypeRefCursor             JDBCType = 2012
	JDBCTypeTimeWithTimeZone      JDBCType = 2013
	JDBCTypeTimestampWithTimezone JDBCType = 2014
)

type MySQLDefCode int64

var (
	COM_BINLOG_DUMP     MySQLDefCode = 18
	COM_CHANGE_USER     MySQLDefCode = 17
	COM_CLOSE_STATEMENT MySQLDefCode = 25
	COM_CONNECT_OUT     MySQLDefCode = 20
	COM_END             MySQLDefCode = 29
	COM_EXECUTE         MySQLDefCode = 23
	COM_FETCH           MySQLDefCode = 28
	COM_LONG_DATA       MySQLDefCode = 24
	COM_PREPARE         MySQLDefCode = 22
	COM_REGISTER_SLAVE  MySQLDefCode = 21
	COM_RESET_STMT      MySQLDefCode = 26
	COM_SET_OPTION      MySQLDefCode = 27
	COM_TABLE_DUMP      MySQLDefCode = 19
	CONNECT             MySQLDefCode = 11
	CREATE_DB           MySQLDefCode = 5 // Not used; deprecated?

	DEBUG          MySQLDefCode = 13
	DELAYED_INSERT MySQLDefCode = 16
	DROP_DB        MySQLDefCode = 6 // Not used; deprecated?

	FIELD_LIST MySQLDefCode = 4 // Not used; deprecated in MySQL 5.7.11 and MySQL 8.0.0.

	FIELD_TYPE_BIT      MySQLDefCode = 16
	FIELD_TYPE_BLOB     MySQLDefCode = 252
	FIELD_TYPE_DATE     MySQLDefCode = 10
	FIELD_TYPE_DATETIME MySQLDefCode = 12

	// Data Types
	FIELD_TYPE_DECIMAL     MySQLDefCode = 0
	FIELD_TYPE_DOUBLE      MySQLDefCode = 5
	FIELD_TYPE_ENUM        MySQLDefCode = 247
	FIELD_TYPE_FLOAT       MySQLDefCode = 4
	FIELD_TYPE_GEOMETRY    MySQLDefCode = 255
	FIELD_TYPE_INT24       MySQLDefCode = 9
	FIELD_TYPE_LONG        MySQLDefCode = 3
	FIELD_TYPE_LONG_BLOB   MySQLDefCode = 251
	FIELD_TYPE_LONGLONG    MySQLDefCode = 8
	FIELD_TYPE_MEDIUM_BLOB MySQLDefCode = 250
	FIELD_TYPE_NEW_DECIMAL MySQLDefCode = 246
	FIELD_TYPE_NEWDATE     MySQLDefCode = 14
	FIELD_TYPE_NULL        MySQLDefCode = 6
	FIELD_TYPE_SET         MySQLDefCode = 248
	FIELD_TYPE_SHORT       MySQLDefCode = 2
	FIELD_TYPE_STRING      MySQLDefCode = 254
	FIELD_TYPE_TIME        MySQLDefCode = 11
	FIELD_TYPE_TIMESTAMP   MySQLDefCode = 7
	FIELD_TYPE_TINY        MySQLDefCode = 1

	// Older data types
	FIELD_TYPE_TINY_BLOB  MySQLDefCode = 249
	FIELD_TYPE_VAR_STRING MySQLDefCode = 253
	FIELD_TYPE_VARCHAR    MySQLDefCode = 15

	// Newer data types
	FIELD_TYPE_YEAR   MySQLDefCode = 13
	FIELD_TYPE_JSON   MySQLDefCode = 245
	INIT_DB           MySQLDefCode = 2
	LENGTH_BLOB       MySQLDefCode = 65535
	LENGTH_LONGBLOB   MySQLDefCode = 4294967295
	LENGTH_MEDIUMBLOB MySQLDefCode = 16777215
	LENGTH_TINYBLOB   MySQLDefCode = 255

	// Limitations
	MAX_ROWS MySQLDefCode = 50000000 // From the MySQL FAQ

	/**
	 * Used to indicate that the server sent no field-level character set information, so the driver should use the connection-level character encoding instead.
	 */
	NO_CHARSET_INFO  MySQLDefCode = -1
	OPEN_CURSOR_FLAG MySQLDefCode = 1
	PING             MySQLDefCode = 14
	PROCESS_INFO     MySQLDefCode = 10 // Not used; deprecated in MySQL 5.7.11 and MySQL 8.0.0.

	PROCESS_KILL MySQLDefCode = 12 // Not used; deprecated in MySQL 5.7.11 and MySQL 8.0.0.

	QUERY  MySQLDefCode = 3
	QUIT   MySQLDefCode = 1
	RELOAD MySQLDefCode = 7 // Not used; deprecated in MySQL 5.7.11 and MySQL 8.0.0.

	SHUTDOWN MySQLDefCode = 8 // Deprecated in MySQL 5.7.9 and MySQL 8.0.0.

	//
	// Constants defined from mysql
	//
	// DB Operations
	SLEEP      MySQLDefCode = 0
	STATISTICS MySQLDefCode = 9
	TIME       MySQLDefCode = 15
)
package datasource

import (
	"database/sql"

	"git.multiverse.io/eventkit/kit/db/datasource/types"
	"github.com/go-sql-driver/mysql"
)

func initDriver() {
	sql.Register(types.ATMySQLDriver, &ATDriver{
		Driver: &Driver{
			DBType:         "mysql",
			Mode:           types.ATMode,
			OriginalDriver: mysql.MySQLDriver{},
		},
	})
}

func init() {
	initDriver()
}

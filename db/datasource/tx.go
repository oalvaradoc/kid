package datasource

import "database/sql/driver"

type Tx struct {
	Conn *Connection
	OriginalTx driver.Tx
}

func (t *Tx) Commit() error {
	return t.OriginalTx.Commit()
}

func (t *Tx) Rollback() error {
	return t.OriginalTx.Rollback()
}
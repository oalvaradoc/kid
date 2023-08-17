package tx

import "git.multiverse.io/eventkit/kit/db/datasource"

type ATTx struct {
	t *datasource.Tx
}


func (a *ATTx) Commit() error {
	return a.t.Commit()
}

func (a *ATTx) Rollback() error {
	return a.t.Rollback()
}
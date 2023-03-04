package transaction_manager

import "github.com/jmoiron/sqlx"

type Transaction interface {
	Commit() error
	Rollback() error
}

type transaction struct {
	tx *sqlx.Tx
}

func NewTransaction(tx *sqlx.Tx) Transaction {
	return transaction{tx: tx}
}

func (t transaction) Commit() error {
	return t.tx.Commit()
}

func (t transaction) Rollback() error {
	return t.tx.Rollback()
}

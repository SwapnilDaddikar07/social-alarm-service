package transaction_manager

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
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

func (t transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

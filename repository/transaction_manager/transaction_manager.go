package transaction_manager

import "github.com/jmoiron/sqlx"

type TransactionManager interface {
	NewTransaction() Transaction
}

type transactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return transactionManager{db: db}
}

func (tm transactionManager) NewTransaction() Transaction {
	return NewTransaction(tm.db.MustBegin())
}

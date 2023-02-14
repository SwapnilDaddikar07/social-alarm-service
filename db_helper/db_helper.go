package db_helper

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

func Connect() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	return sqlx.Connect("mysql", dsn)
}

package db_helper

import (
	"github.com/jmoiron/sqlx"
)

// Connect TODO Read username , password , HOST from ENV variables
func Connect() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", "alarm_user:alarm_pswd@localhost:3306/alarm_database")
}

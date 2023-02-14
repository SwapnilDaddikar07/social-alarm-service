package db_model

import (
	"database/sql"
)

type PublicNonExpiredAlarms struct {
	AlarmId       string       `db:"alarm_id"`
	Description   string       `db:"description"`
	StartDateTime sql.NullTime `db:"alarm_start_datetime"`
	Mon           int          `db:"mon"`
	Tue           int          `db:"tue"`
	Wed           int          `db:"wed"`
	Thu           int          `db:"thu"`
	Fri           int          `db:"fri"`
	Sat           int          `db:"sat"`
	Sun           int          `db:"sun"`
}

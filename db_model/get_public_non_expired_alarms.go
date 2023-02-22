package db_model

import (
	"database/sql"
)

type PublicNonExpiredRepeatingAlarms struct {
	AlarmId          string       `db:"alarm_id"`
	Description      string       `db:"description"`
	StartDateTime    sql.NullTime `db:"alarm_start_datetime"`
	MonSystemAlarmId int          `db:"mon_system_alarm_id"`
	TueSystemAlarmId int          `db:"tue_system_alarm_id"`
	WedSystemAlarmId int          `db:"wed_system_alarm_id"`
	ThuSystemAlarmId int          `db:"thu_system_alarm_id"`
	FriSystemAlarmId int          `db:"fri_system_alarm_id"`
	SatSystemAlarmId int          `db:"sat_system_alarm_id"`
	SunSystemAlarmId int          `db:"sun_system_alarm_id"`
}

type PublicNonExpiredNonRepeatingAlarms struct {
	AlarmId       string       `db:"alarm_id"`
	Description   string       `db:"description"`
	StartDateTime sql.NullTime `db:"alarm_start_datetime"`
	SystemAlarmId int          `db:"system_alarm_id"`
}

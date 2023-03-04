package db_model

import (
	"database/sql"
)

type PublicNonExpiredRepeatingAlarms struct {
	AlarmId          string       `db:"alarm_id"`
	Description      string       `db:"description"`
	StartDateTime    sql.NullTime `db:"alarm_start_datetime"`
	MonDeviceAlarmId int          `db:"mon_device_alarm_id"`
	TueDeviceAlarmId int          `db:"tue_device_alarm_id"`
	WedDeviceAlarmId int          `db:"wed_device_alarm_id"`
	ThuDeviceAlarmId int          `db:"thu_device_alarm_id"`
	FriDeviceAlarmId int          `db:"fri_device_alarm_id"`
	SatDeviceAlarmId int          `db:"sat_device_alarm_id"`
	SunDeviceAlarmId int          `db:"sun_device_alarm_id"`
}

type PublicNonExpiredNonRepeatingAlarms struct {
	AlarmId       string       `db:"alarm_id"`
	Description   string       `db:"description"`
	StartDateTime sql.NullTime `db:"alarm_start_datetime"`
}

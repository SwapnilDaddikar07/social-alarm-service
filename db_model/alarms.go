package db_model

import (
	"database/sql"
	"social-alarm-service/constants"
)

type Alarms struct {
	AlarmID            string                    `db:"alarm_id"`
	UserID             string                    `db:"user_id"`
	Visibility         constants.AlarmVisibility `db:"visibility"`
	Description        string                    `db:"description"`
	Status             string                    `db:"status"`
	AlarmStartDateTime sql.NullTime              `db:"alarm_start_datetime"`
	CreatedAt          sql.NullTime              `db:"created_at"`
	MonDeviceAlarmId   int                       `db:"mon_device_alarm_id"`
	TueDeviceAlarmId   int                       `db:"tue_device_alarm_id"`
	WedDeviceAlarmId   int                       `db:"wed_device_alarm_id"`
	ThuDeviceAlarmId   int                       `db:"thu_device_alarm_id"`
	FriDeviceAlarmId   int                       `db:"fri_device_alarm_id"`
	SatDeviceAlarmId   int                       `db:"sat_device_alarm_id"`
	SunDeviceAlarmId   int                       `db:"sun_device_alarm_id"`
}

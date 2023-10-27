package response_model

import (
	"social-alarm-service/constants"
	"social-alarm-service/db_model"
	"time"
)

type NonRepeatingAlarms struct {
	AlarmID                   string                    `json:"alarm_id"`
	UserID                    string                    `json:"user_id"`
	Visibility                constants.AlarmVisibility `json:"visibility"`
	Description               string                    `json:"description"`
	Status                    string                    `json:"status"`
	AlarmStartDateTime        time.Time                 `json:"alarm_start_datetime"`
	CreatedAt                 time.Time                 `json:"created_at"`
	NonRepeatingDeviceAlarmId int                       `json:"device_alarm_id"`
}

type RepeatingAlarms struct {
	AlarmID            string                    `json:"alarm_id"`
	UserID             string                    `json:"user_id"`
	Visibility         constants.AlarmVisibility `json:"visibility"`
	Description        string                    `json:"description"`
	Status             string                    `json:"status"`
	AlarmStartDateTime time.Time                 `json:"alarm_start_datetime"`
	CreatedAt          time.Time                 `json:"created_at"`
	MonDeviceAlarmId   *int                      `json:"mon_device_alarm_id,omitempty"`
	TueDeviceAlarmId   *int                      `json:"tue_device_alarm_id,omitempty"`
	WedDeviceAlarmId   *int                      `json:"wed_device_alarm_id,omitempty"`
	ThuDeviceAlarmId   *int                      `json:"thu_device_alarm_id,omitempty"`
	FriDeviceAlarmId   *int                      `json:"fri_device_alarm_id,omitempty"`
	SatDeviceAlarmId   *int                      `json:"sat_device_alarm_id,omitempty"`
	SunDeviceAlarmId   *int                      `json:"sun_device_alarm_id,omitempty"`
}

type GetAllAlarms struct {
	RepeatingAlarms    []RepeatingAlarms    `json:"repeating_alarms"`
	NonRepeatingAlarms []NonRepeatingAlarms `json:"non_repeating_alarms"`
}

func MapToRepeatingAlarms(dbAlarms []db_model.Alarms) []RepeatingAlarms {
	ra := make([]RepeatingAlarms, 0)
	for _, dbAlarm := range dbAlarms {
		ra = append(ra, MapToRepeatingAlarm(dbAlarm))
	}
	return ra
}

func MapToRepeatingAlarm(dbAlarms db_model.Alarms) RepeatingAlarms {
	ra := RepeatingAlarms{}
	ra.AlarmID = dbAlarms.AlarmID
	ra.UserID = dbAlarms.UserID
	ra.Description = dbAlarms.Description
	ra.Visibility = dbAlarms.Visibility
	ra.Status = dbAlarms.Status
	ra.CreatedAt = dbAlarms.CreatedAt.Time
	ra.AlarmStartDateTime = dbAlarms.AlarmStartDateTime.Time
	ra.MonDeviceAlarmId = &dbAlarms.MonDeviceAlarmId
	ra.TueDeviceAlarmId = &dbAlarms.TueDeviceAlarmId
	ra.WedDeviceAlarmId = &dbAlarms.WedDeviceAlarmId
	ra.ThuDeviceAlarmId = &dbAlarms.ThuDeviceAlarmId
	ra.FriDeviceAlarmId = &dbAlarms.FriDeviceAlarmId
	ra.SatDeviceAlarmId = &dbAlarms.SatDeviceAlarmId
	ra.SunDeviceAlarmId = &dbAlarms.SunDeviceAlarmId
	return ra
}

func MapToNonRepeatingAlarms(dbAlarms []db_model.Alarms) []NonRepeatingAlarms {
	nra := make([]NonRepeatingAlarms, 0)
	for _, dbAlarm := range dbAlarms {
		nra = append(nra, MapToNonRepeatingAlarm(dbAlarm))
	}
	return nra
}

func MapToNonRepeatingAlarm(dbAlarms db_model.Alarms) NonRepeatingAlarms {
	nra := NonRepeatingAlarms{}
	nra.AlarmID = dbAlarms.AlarmID
	nra.UserID = dbAlarms.UserID
	nra.Description = dbAlarms.Description
	nra.Visibility = dbAlarms.Visibility
	nra.Status = dbAlarms.Status
	nra.CreatedAt = dbAlarms.CreatedAt.Time
	nra.NonRepeatingDeviceAlarmId = dbAlarms.NonRepeatingDeviceAlarmId
	return nra
}

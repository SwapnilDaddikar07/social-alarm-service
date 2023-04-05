package response_model

import (
	"social-alarm-service/db_model"
	"time"
)

type EligibleAlarmsResponse struct {
	UserId                     string                       `json:"user_id"`
	EligibleRepeatingAlarms    []EligibleRepeatingAlarms    `json:"repeating_alarms"`
	EligibleNonRepeatingAlarms []EligibleNonRepeatingAlarms `json:"non_repeating_alarms"`
}

type EligibleRepeatingAlarms struct {
	AlarmId            string    `json:"alarm_id"`
	AlarmStartDateTime time.Time `json:"alarm_start_date_time"`
	Description        string    `json:"description"`
	Schedules          []string  `json:"schedules"`
}

type EligibleNonRepeatingAlarms struct {
	AlarmId            string    `json:"alarm_id"`
	AlarmStartDateTime time.Time `json:"alarm_start_date_time"`
	Description        string    `json:"description"`
}

func MapToEligibleNonRepeatingAlarms(alarms []db_model.Alarms) []EligibleNonRepeatingAlarms {
	eligibleNonRepeatingAlarms := make([]EligibleNonRepeatingAlarms, 0)
	for _, alarm := range alarms {
		eligibleNonRepeatingAlarms = append(eligibleNonRepeatingAlarms, MapToEligibleNonRepeatingAlarm(alarm))
	}
	return eligibleNonRepeatingAlarms
}

func MapToEligibleNonRepeatingAlarm(alarm db_model.Alarms) EligibleNonRepeatingAlarms {
	return EligibleNonRepeatingAlarms{
		AlarmId:            alarm.AlarmID,
		AlarmStartDateTime: alarm.AlarmStartDateTime.Time,
		Description:        alarm.Description,
	}
}

func MapToEligibleRepeatingAlarms(alarms []db_model.Alarms) []EligibleRepeatingAlarms {
	eligibleNonRepeatingAlarms := make([]EligibleRepeatingAlarms, 0)
	for _, alarm := range alarms {
		eligibleNonRepeatingAlarms = append(eligibleNonRepeatingAlarms, MapToEligibleRepeatingAlarm(alarm))
	}
	return eligibleNonRepeatingAlarms
}

func MapToEligibleRepeatingAlarm(alarm db_model.Alarms) EligibleRepeatingAlarms {
	return EligibleRepeatingAlarms{
		AlarmId:            alarm.AlarmID,
		AlarmStartDateTime: alarm.AlarmStartDateTime.Time,
		Description:        alarm.Description,
		Schedules:          generateSchedules(alarm),
	}
}

func generateSchedules(alarm db_model.Alarms) []string {
	responseSchedule := make([]string, 0)

	if alarm.MonDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Mon")
	}
	if alarm.TueDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Tue")
	}
	if alarm.WedDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Wed")
	}
	if alarm.ThuDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Thu")
	}
	if alarm.FriDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Fri")
	}
	if alarm.SatDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Sat")
	}
	if alarm.SunDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Sun")
	}

	return responseSchedule
}

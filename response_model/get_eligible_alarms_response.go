package response_model

import (
	"social-alarm-service/db_model"
	"time"
)

type EligibleAlarmsResponse struct {
	AlarmId            string    `json:"alarm_id"`
	AlarmStartDateTime time.Time `json:"alarm_start_date_time"`
	Description        string    `json:"description"`
	Schedules          []string  `json:"schedules"`
}

func MapNonRepeatingAlarmsToEligibleAlarmsResponseList(publicNonExpiredNonRepeatingAlarms []db_model.PublicNonExpiredNonRepeatingAlarms) []EligibleAlarmsResponse {
	eligibleAlarmsResponse := make([]EligibleAlarmsResponse, 0)

	for _, entry := range publicNonExpiredNonRepeatingAlarms {
		eligibleAlarmsResponse = append(eligibleAlarmsResponse, EligibleAlarmsResponse{
			AlarmId:            entry.AlarmId,
			AlarmStartDateTime: entry.StartDateTime.Time,
			Description:        entry.Description,
			Schedules:          nil,
		})
	}
	return eligibleAlarmsResponse
}

func MapRepeatingAlarmsToEligibleAlarmsResponseList(publicNonExpiredRepeatingAlarms []db_model.PublicNonExpiredRepeatingAlarms) []EligibleAlarmsResponse {
	eligibleAlarmsResponse := make([]EligibleAlarmsResponse, 0)

	for _, entry := range publicNonExpiredRepeatingAlarms {
		eligibleAlarmsResponse = append(eligibleAlarmsResponse, MapToEligibleRepeatingAlarmsResponse(entry))
	}
	return eligibleAlarmsResponse
}

func MapToEligibleRepeatingAlarmsResponse(publicNonExpiredAlarms db_model.PublicNonExpiredRepeatingAlarms) EligibleAlarmsResponse {
	return EligibleAlarmsResponse{
		AlarmId:            publicNonExpiredAlarms.AlarmId,
		AlarmStartDateTime: publicNonExpiredAlarms.StartDateTime.Time,
		Description:        publicNonExpiredAlarms.Description,
		Schedules:          generateSchedules(publicNonExpiredAlarms),
	}
}

func generateSchedules(dbSchedules db_model.PublicNonExpiredRepeatingAlarms) []string {
	responseSchedule := make([]string, 0)

	if dbSchedules.MonDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Mon")
	}
	if dbSchedules.TueDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Tue")
	}
	if dbSchedules.WedDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Wed")
	}
	if dbSchedules.ThuDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Thu")
	}
	if dbSchedules.FriDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Fri")
	}
	if dbSchedules.SatDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Sat")
	}
	if dbSchedules.SunDeviceAlarmId >= 0 {
		responseSchedule = append(responseSchedule, "Sun")
	}

	return responseSchedule
}

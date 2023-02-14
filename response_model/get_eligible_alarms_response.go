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

func MapToEligibleAlarmsResponseList(publicNonExpiredAlarms []db_model.PublicNonExpiredAlarms) []EligibleAlarmsResponse {
	eligibleAlarmsResponse := make([]EligibleAlarmsResponse, 0)

	for _, entry := range publicNonExpiredAlarms {
		eligibleAlarmsResponse = append(eligibleAlarmsResponse, MapToEligibleAlarmsResponse(entry))
	}
	return eligibleAlarmsResponse
}

func MapToEligibleAlarmsResponse(publicNonExpiredAlarms db_model.PublicNonExpiredAlarms) EligibleAlarmsResponse {
	return EligibleAlarmsResponse{
		AlarmId:            publicNonExpiredAlarms.AlarmId,
		AlarmStartDateTime: publicNonExpiredAlarms.StartDateTime.Time,
		Description:        publicNonExpiredAlarms.Description,
		Schedules:          nil,
	}
}

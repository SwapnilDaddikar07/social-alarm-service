package model

import "time"

type AlarmResponse struct {
	AlarmId            string    `json:"alarm_id"`
	AlarmStartDateTime time.Time `json:"alarm_start_date_time"`
	Description        string    `json:"description"`
	Schedules          []string  `json:"schedules"`
}

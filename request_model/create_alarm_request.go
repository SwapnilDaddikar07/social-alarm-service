package request_model

type CreateAlarmRequest struct {
	AlarmStartDateTime        string `json:"alarm_start_date_time" binding:"required"`
	Visibility                bool   `json:"visibility"`
	Description               string `json:"description"`
	UserId                    string `json:"user_id" binding:"required"` // change this after moving ID to token.
	RepeatingSystemAlarmIds   []int  `json:"repeating_system_alarm_ids"`
	NonRepeatingSystemAlarmId int    `json:"non_repeating_system_alarm_id"`
}

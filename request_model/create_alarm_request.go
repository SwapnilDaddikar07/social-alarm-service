package request_model

type CreateAlarmRequest struct {
	AlarmStartDateTime        string `json:"alarm_start_date_time"`
	Visibility                bool   `json:"visibility"`
	Description               string `json:"description"`
	UserId                    string `json:"user_id"`
	RepeatingSystemAlarmIds   []int  `json:"repeating_system_alarm_ids"`
	NonRepeatingSystemAlarmId int    `json:"non_repeating_system_alarm_id"`
}

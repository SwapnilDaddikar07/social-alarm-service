package request_model

type UpdateAlarmStatus struct {
	UserId  string `json:"user_id" binding:"required"`
	AlarmId string `json:"alarm_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

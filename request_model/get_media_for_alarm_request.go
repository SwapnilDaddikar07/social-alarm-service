package request_model

type GetMediaForAlarm struct {
	AlarmId string `json:"alarm_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

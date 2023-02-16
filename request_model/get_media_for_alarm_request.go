package request_model

type GetMediaForAlarm struct {
	AlarmId string `json:"alarm_id" binding:"required"`
}

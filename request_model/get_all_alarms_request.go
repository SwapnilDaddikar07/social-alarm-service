package request_model

type GetAllAlarmsRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

package request_model

type GetEligibleAlarmsRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

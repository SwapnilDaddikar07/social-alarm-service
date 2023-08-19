package request_model

type GetProfilesRequest struct {
	UserId       string   `json:"user_id" binding:"required"`
	PhoneNumbers []string `json:"phone_numbers" binding:"required"`
}

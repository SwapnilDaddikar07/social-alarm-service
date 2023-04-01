package request_model

type GetProfilesRequest struct {
	PhoneNumbers []string `json:"phone_numbers" binding:"required"`
}

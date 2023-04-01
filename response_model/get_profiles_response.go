package response_model

import "social-alarm-service/db_model"

type GetProfilesResponse []Profile

type Profile struct {
	DisplayName string `json:"display_name"`
	UserId      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}

func MapToProfilesResponse(users []db_model.User) []Profile {
	profilesResponse := make([]Profile, 0)
	for _, user := range users {
		profilesResponse = append(profilesResponse, MapToProfile(user))
	}
	return profilesResponse
}

func MapToProfile(user db_model.User) Profile {
	return Profile{
		DisplayName: user.DisplayName,
		UserId:      user.UserId,
		PhoneNumber: user.PhoneNumber,
	}
}

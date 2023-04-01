package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/response_model"
)

type UserService interface {
	GetProfiles(ctx *gin.Context, phoneNumbers []string) (response_model.GetProfilesResponse, *error2.ASError)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (us userService) GetProfiles(ctx *gin.Context, phoneNumbers []string) (response_model.GetProfilesResponse, *error2.ASError) {
	profiles, err := us.userRepo.GetProfiles(ctx, phoneNumbers)
	if err != nil {
		fmt.Printf("repo error when fetching profile details %v", err)
		return response_model.GetProfilesResponse{}, error2.InternalServerError("unable to fetch profiles")
	}
	fmt.Printf("received %d number and retrived %d profiles", len(phoneNumbers), len(profiles))

	return response_model.MapToProfilesResponse(profiles), nil
}

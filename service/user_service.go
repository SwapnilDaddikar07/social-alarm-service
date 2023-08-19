package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/response_model"
	"social-alarm-service/utils"
	"strconv"
)

type UserService interface {
	GetProfiles(ctx *gin.Context, userId string, phoneNumbers []string) (response_model.GetProfilesResponse, *error2.ASError)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

// GetProfiles TODO can only be hit if user is logged in.
func (us userService) GetProfiles(ctx *gin.Context, userId string, phoneNumbers []string) (response_model.GetProfilesResponse, *error2.ASError) {
	user, dbErr := us.userRepo.GetUser(ctx, userId)
	if dbErr != nil {
		return response_model.GetProfilesResponse{}, error2.InternalServerError("unable to fetch user profile for given user id")
	}

	regionCode, err := utils.GetRegionCode(user.PhoneNumber)
	if err != nil {
		return response_model.GetProfilesResponse{}, error2.InternalServerError("unable to fetch profiles")
	}

	formattedNumbers := us.formatNumbers(regionCode, phoneNumbers)

	profiles, dbErr := us.userRepo.GetProfiles(ctx, formattedNumbers)
	if dbErr != nil {
		fmt.Printf("repo error when fetching profile details %v", dbErr)
		return response_model.GetProfilesResponse{}, error2.InternalServerError("unable to fetch profiles")
	}

	fmt.Printf("received %d number and retrived %d profiles", len(phoneNumbers), len(profiles))
	return response_model.MapToProfilesResponse(profiles), nil
}

func (us userService) formatNumbers(defaultRegionCode string, nonformattedMobileNumbers []string) []string {
	formattedMobileNumbers := make([]string, 0)

	for _, mobileNum := range nonformattedMobileNumbers {
		p, err := phonenumbers.Parse(mobileNum, defaultRegionCode)
		if err != nil {
			fmt.Println(fmt.Sprintf("could not parse the number %s ", mobileNum))
			continue
		}

		formattedNumber := fmt.Sprintf("%s%s%s", "+", fmt.Sprint(p.GetCountryCode()), strconv.FormatUint(p.GetNationalNumber(), 10))
		fmt.Println(fmt.Sprintf("unformatted number is %s and formatted number is %s", mobileNum, formattedNumber))

		formattedMobileNumbers = append(formattedMobileNumbers, formattedNumber)
	}

	return formattedMobileNumbers
}

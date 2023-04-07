package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	error2 "social-alarm-service/error"
	"social-alarm-service/request_model"
	"social-alarm-service/service"
)

type UserController interface {
	GetProfiles(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return userController{userService: userService}
}

func (uc userController) GetProfiles(ctx *gin.Context) {
	request := request_model.GetProfilesRequest{}

	bindingErr := ctx.ShouldBindWith(&request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if len(request.PhoneNumbers) == 0 {
		fmt.Println("no mobile numbers present")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, error2.NoPhoneNumbersInRequest)
		return
	}

	userProfiles, serviceErr := uc.userService.GetProfiles(ctx, request.PhoneNumbers)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, userProfiles)
}

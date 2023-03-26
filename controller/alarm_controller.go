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

type AlarmController interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context)
	CreateAlarm(ctx *gin.Context)
	UpdateAlarmStatus(ctx *gin.Context)
	GetAllAlarms(ctx *gin.Context)
}

type alarmController struct {
	alarmService service.AlarmService
}

func NewAlarmController(alarmService service.AlarmService) AlarmController {
	return alarmController{alarmService: alarmService}
}

func (ac alarmController) GetPublicNonExpiredAlarms(ctx *gin.Context) {
	request := request_model.GetEligibleAlarmsRequest{}

	bindingErr := ctx.ShouldBindBodyWith(&request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	allEligibleAlarms, serviceErr := ac.alarmService.GetPublicNonExpiredAlarms(ctx, request.UserId)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, allEligibleAlarms)
}

func (ac alarmController) CreateAlarm(ctx *gin.Context) {
	request := &request_model.CreateAlarmRequest{}

	bindingErr := ctx.ShouldBindWith(request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, error2.BadRequestError("invalid request"))
		return
	}

	createAlarmResponse, serviceErr := ac.alarmService.CreateAlarm(ctx, *request)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(http.StatusCreated, createAlarmResponse)
}

func (ac alarmController) UpdateAlarmStatus(ctx *gin.Context) {
	request := &request_model.UpdateAlarmStatus{}

	bindingErr := ctx.ShouldBindWith(request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, error2.BadRequestError("invalid request"))
		return
	}

	serviceErr := ac.alarmService.UpdateStatus(ctx, request.AlarmId, request.UserId, request.Status)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func (ac alarmController) GetAllAlarms(ctx *gin.Context) {
	request := &request_model.GetAllAlarmsRequest{}

	bindingErr := ctx.ShouldBindWith(request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, error2.BadRequestError("invalid request"))
		return
	}

	allAlarms, serviceErr := ac.alarmService.GetAllAlarms(ctx, request.UserId)
	if serviceErr != nil {
		fmt.Printf("service error %v\n", serviceErr)
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, allAlarms)
}

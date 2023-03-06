package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	error2 "social-alarm-service/error"
	"social-alarm-service/request_model"
	"social-alarm-service/service"
)

type AlarmController interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context)
	CreateAlarm(ctx *gin.Context)
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
		ctx.AbortWithStatus(400)
		return
	}

	allEligibleAlarms, serviceErr := ac.alarmService.GetPublicNonExpiredAlarms(ctx, request.UserId)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(200, allEligibleAlarms)
}

func (ac alarmController) CreateAlarm(ctx *gin.Context) {
	request := &request_model.CreateAlarmRequest{}

	bindingErr := ctx.ShouldBindWith(request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatusJSON(400, error2.BadRequestError("invalid request"))
		return
	}

	createAlarmResponse, serviceErr := ac.alarmService.CreateAlarm(ctx, *request)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(201, createAlarmResponse)

}

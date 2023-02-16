package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"social-alarm-service/request_model"
	"social-alarm-service/service"
)

type AlarmController interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context)
	GetMediaForAlarm(ctx *gin.Context)
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

func (ac alarmController) GetMediaForAlarm(ctx *gin.Context) {
	request := request_model.GetMediaForAlarm{}

	bindingErr := ctx.ShouldBindWith(&request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatus(400)
		return
	}

	alarmMedia, serviceErr := ac.alarmService.GetMediaForAlarm(ctx, request.AlarmId)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(200, alarmMedia)
}

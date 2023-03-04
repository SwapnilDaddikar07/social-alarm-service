package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"social-alarm-service/request_model"
	"social-alarm-service/service"
)

type AlarmMediaController interface {
	GetMediaForAlarm(ctx *gin.Context)
}

type alarmMediaController struct {
	service service.AlarmMediaService
}

func NewAlarmMediaController(service service.AlarmMediaService) AlarmMediaController {
	return alarmMediaController{service: service}
}

func (amc alarmMediaController) GetMediaForAlarm(ctx *gin.Context) {
	request := request_model.GetMediaForAlarm{}

	bindingErr := ctx.ShouldBindWith(&request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatus(400)
		return
	}

	alarmMedia, serviceErr := amc.service.GetMediaForAlarm(ctx, request.AlarmId)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(200, alarmMedia)
}

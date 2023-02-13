package controller

import (
	"github.com/gin-gonic/gin"
	"social-alarm-service/service"
)

type AlarmController interface {
	GetAllAlarms(ctx *gin.Context)
}

type alarmController struct {
	alarmService service.AlarmService
}

func NewAlarmController(alarmService service.AlarmService) AlarmController {
	return alarmController{alarmService: alarmService}
}

func (ac alarmController) GetAllAlarms(ctx *gin.Context) {

}

package controller

import (
	"github.com/gin-gonic/gin"
	"social-alarm-service/service"
)

type AlarmController interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context)
}

type alarmController struct {
	alarmService service.AlarmService
}

func NewAlarmController(alarmService service.AlarmService) AlarmController {
	return alarmController{alarmService: alarmService}
}

func (ac alarmController) GetPublicNonExpiredAlarms(ctx *gin.Context) {

}

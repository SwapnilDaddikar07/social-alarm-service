package service

import (
	"github.com/gin-gonic/gin"
	"social-alarm-service/repository"
	"social-alarm-service/response_model"
)

type AlarmService interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]response_model.EligibleAlarmsResponse, error)
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, error)
}

type alarmService struct {
	alarmRepository repository.AlarmRepository
}

func NewAlarmService(alarmRepository repository.AlarmRepository) AlarmService {
	return alarmService{alarmRepository: alarmRepository}
}

func (as alarmService) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]response_model.EligibleAlarmsResponse, error) {
	allPublicNonExpiredAlarms, err := as.alarmRepository.GetPublicNonExpiredAlarms(ctx, userId)
	if err != nil {
		return []response_model.EligibleAlarmsResponse{}, err
	}
	return response_model.MapToEligibleAlarmsResponseList(allPublicNonExpiredAlarms), nil
}

func (as alarmService) GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, error) {
	alarmMedia, err := as.alarmRepository.GetMediaForAlarm(ctx, alarmId)
	if err != nil {
		return []response_model.MediaForAlarm{}, err
	}
	return response_model.MapToMediaForAlarmResponseList(alarmMedia), nil
}

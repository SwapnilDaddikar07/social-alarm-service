package service

import (
	"github.com/gin-gonic/gin"
	"social-alarm-service/model"
	"social-alarm-service/repository"
)

type AlarmService interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error)
}

type alarmService struct {
	alarmRepository repository.AlarmRepository
}

func NewAlarmService(alarmRepository repository.AlarmRepository) AlarmService {
	return alarmService{alarmRepository: alarmRepository}
}

func (as alarmService) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error) {
	dbAlarmResponse, err := as.alarmRepository.GetPublicNonExpiredAlarms(ctx, userId)
	if err != nil {
		return []model.AlarmResponse{}, err
	}
	return dbAlarmResponse, nil
}

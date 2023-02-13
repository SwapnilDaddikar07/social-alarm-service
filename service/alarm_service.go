package service

import (
	"github.com/gin-gonic/gin"
	"social-alarm-service/model"
	"social-alarm-service/repository"
)

type AlarmService interface {
	GetAllAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error)
}

type alarmService struct {
	alarmRepository repository.AlarmRepository
}

func NewAlarmService(alarmRepository repository.AlarmRepository) AlarmService {
	return alarmService{alarmRepository: alarmRepository}
}

func (as alarmService) GetAllAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error) {
	dbAlarmResponse, err := as.alarmRepository.GetAllAlarms(ctx, userId)
	if err != nil {
		return []model.AlarmResponse{}, err
	}
	return dbAlarmResponse, nil
}

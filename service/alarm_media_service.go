package service

import (
	"github.com/gin-gonic/gin"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/response_model"
)

type AlarmMediaService interface {
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError)
}

type alarmMediaService struct {
	alarmMediaRepo repository.AlarmMediaRepository
}

func NewAlarmMediaService(alarmMediaRepo repository.AlarmMediaRepository) AlarmMediaService {
	return alarmMediaService{alarmMediaRepo: alarmMediaRepo}
}

func (as alarmMediaService) GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError) {
	alarmMedia, err := as.alarmMediaRepo.GetMediaForAlarm(ctx, alarmId)
	if err != nil {
		return []response_model.MediaForAlarm{}, error2.InternalServerError("db fetch error when getting all media associated with given alarm id")
	}
	return response_model.MapToMediaForAlarmResponseList(alarmMedia), nil
}

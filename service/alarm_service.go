package service

import (
	"github.com/gin-gonic/gin"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/request_model"
	"social-alarm-service/response_model"
	"time"
)

type AlarmService interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]response_model.EligibleAlarmsResponse, *error2.ASError)
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError)
	CreateAlarm(ctx *gin.Context, request request_model.CreateAlarmRequest) (response_model.CreateAlarmResponse, *error2.ASError)
}

type alarmService struct {
	alarmRepository repository.AlarmRepository
}

func NewAlarmService(alarmRepository repository.AlarmRepository) AlarmService {
	return alarmService{alarmRepository: alarmRepository}
}

func (as alarmService) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]response_model.EligibleAlarmsResponse, *error2.ASError) {
	publicNonExpiredRepeatingAlarms, publicNonExpiredNonRepeatingAlarms, err := as.alarmRepository.GetPublicNonExpiredAlarms(ctx, userId)
	if err != nil {
		return []response_model.EligibleAlarmsResponse{}, error2.InternalServerError("db fetch error when getting public non expired alarms for given user id")
	}

	eligibleAlarms := response_model.MapRepeatingAlarmsToEligibleAlarmsResponseList(publicNonExpiredRepeatingAlarms)
	eligibleAlarms = append(eligibleAlarms, response_model.MapNonRepeatingAlarmsToEligibleAlarmsResponseList(publicNonExpiredNonRepeatingAlarms)...)

	return eligibleAlarms, nil
}

func (as alarmService) GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError) {
	alarmMedia, err := as.alarmRepository.GetMediaForAlarm(ctx, alarmId)
	if err != nil {
		return []response_model.MediaForAlarm{}, error2.InternalServerError("db fetch error when getting all media associated with given alarm id")
	}
	return response_model.MapToMediaForAlarmResponseList(alarmMedia), nil
}

func (as alarmService) CreateAlarm(ctx *gin.Context, request request_model.CreateAlarmRequest) (response_model.CreateAlarmResponse, *error2.ASError) {
	if len(request.RepeatingSystemAlarmIds) == 0 && request.NonRepeatingSystemAlarmId == nil {
		return response_model.CreateAlarmResponse{}, error2.AlarmIdMissing
	}
	if len(request.RepeatingSystemAlarmIds) >= 1 && (request.NonRepeatingSystemAlarmId != nil) {
		return response_model.CreateAlarmResponse{}, error2.InvalidAlarmTypeError
	}
	//check the DB len
	if len(request.Description) > 50 {
		return response_model.CreateAlarmResponse{}, error2.DescriptionTooLongError
	}
	userExists, dbError := as.alarmRepository.UserExists(ctx, request.UserId)
	if dbError != nil {
		return response_model.CreateAlarmResponse{}, error2.InternalServerError("db fetch error")
	}
	if !userExists {
		return response_model.CreateAlarmResponse{}, error2.InvalidUserIdError
	}

	if len(request.RepeatingSystemAlarmIds) > 0 {
		saveRepeatingAlarm(ctx, request.RepeatingSystemAlarmIds, request.UserId, request.AlarmStartDateTime)
	} else {
		saveNonRepeatingAlarm(ctx, *request.NonRepeatingSystemAlarmId, request.UserId, request.AlarmStartDateTime)
	}
	panic("")
}

func saveNonRepeatingAlarm(ctx *gin.Context, systemAlarmId int, userId string, alarmStartDateTime time.Time) {

}

func saveRepeatingAlarm(ctx *gin.Context, repeatingSystemAlarmIds []int, userId string, alarmStartDateTime time.Time) {

}

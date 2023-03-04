package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/repository/transaction_manager"
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
	alarmRepository    repository.AlarmRepository
	transactionManager transaction_manager.TransactionManager
}

func NewAlarmService(alarmRepository repository.AlarmRepository, transactionManager transaction_manager.TransactionManager) AlarmService {
	return alarmService{alarmRepository: alarmRepository, transactionManager: transactionManager}
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
	if !request.ContainsAtleastOneRepeatingAlarm() && request.NonRepeatingSystemAlarmId == nil {
		return response_model.CreateAlarmResponse{}, error2.AlarmIdMissing
	}

	if request.ContainsAtleastOneRepeatingAlarm() && (request.NonRepeatingSystemAlarmId != nil) {
		return response_model.CreateAlarmResponse{}, error2.InvalidAlarmTypeError
	}

	//TODO check the DB len
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

	//TODO Add the layout
	parsedTime, parseErr := time.Parse("", request.AlarmStartDateTime)
	if parseErr != nil {
		return response_model.CreateAlarmResponse{}, error2.InvalidAlarmDateTimeFormat
	}

	isPrivateAlarm := "F"
	if request.Private {
		isPrivateAlarm = "T"
	}

	//TODO move this to UTIL else code becomes untestable.
	alarmID := uuid.New().String()

	transaction := as.transactionManager.NewTransaction()

	createAlarmDBError := as.alarmRepository.CreateAlarmMetadata(ctx, transaction, alarmID, request.UserId, parsedTime, isPrivateAlarm, request.Description)
	if createAlarmDBError != nil {
		transaction.Rollback()
		return response_model.CreateAlarmResponse{}, error2.InternalServerError("error creating alarm")
	}

	var deviceAlarmSaveError error
	if request.ContainsAtleastOneRepeatingAlarm() {
		deviceAlarmSaveError = as.saveRepeatingDeviceAlarmIds(ctx, transaction, request.RepeatingSystemAlarmIds, alarmID)
	} else {
		deviceAlarmSaveError = as.saveNonRepeatingDeviceAlarmId(ctx, transaction, *request.NonRepeatingSystemAlarmId, alarmID)
	}

	if deviceAlarmSaveError != nil {
		return response_model.CreateAlarmResponse{}, error2.InternalServerError("could not save alarm.")
		transaction.Rollback()
	}
	commitError := transaction.Commit()
	if commitError != nil {
		return response_model.CreateAlarmResponse{}, error2.InternalServerError("db commit failed.")
	}

	fmt.Println("alarm saved successfully.")
	return response_model.CreateAlarmResponse{AlarmId: alarmID}, nil
}

func (as alarmService) saveRepeatingDeviceAlarmIds(ctx *gin.Context, transaction transaction_manager.Transaction, repeatingAlarmIds request_model.RepeatingSystemAlarmIds, alarmId string) error {
	return as.alarmRepository.InsertRepeatingDeviceAlarmIDs(ctx, transaction, alarmId, repeatingAlarmIds.MapToDBModel())
}

func (as alarmService) saveNonRepeatingDeviceAlarmId(ctx *gin.Context, transaction transaction_manager.Transaction, deviceAlarmId int, alarmId string) error {
	return as.alarmRepository.InsertNonRepeatingDeviceAlarmID(ctx, transaction, alarmId, deviceAlarmId)
}

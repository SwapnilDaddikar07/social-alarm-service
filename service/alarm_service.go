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

func (as alarmService) CreateAlarm(ctx *gin.Context, request request_model.CreateAlarmRequest) (response response_model.CreateAlarmResponse, asError *error2.ASError) {
	asError = as.validateCreateAlarmRequest(request)
	if asError != nil {
		return
	}

	userExists, dbError := as.alarmRepository.UserExists(ctx, request.UserId)
	if dbError != nil {
		asError = error2.InternalServerError("db fetch error")
		return
	}
	if !userExists {
		asError = error2.InvalidUserIdError
		return
	}

	alarmId, asError := as.saveAlarm(ctx, request)
	if asError != nil {
		return
	}

	fmt.Println("alarm saved successfully.")
	return response_model.CreateAlarmResponse{AlarmId: alarmId}, nil
}

func (as alarmService) validateCreateAlarmRequest(request request_model.CreateAlarmRequest) *error2.ASError {
	if !request.RepeatingDeviceAlarmIds.ContainsAtleastOneRepeatingAlarm() && request.NonRepeatingDeviceAlarmId == nil {
		return error2.AlarmIdMissing
	}
	if request.RepeatingDeviceAlarmIds.ContainsAtleastOneRepeatingAlarm() && (request.NonRepeatingDeviceAlarmId != nil) {
		return error2.InvalidAlarmTypeError
	}
	//TODO check the DB len
	if len(request.Description) > 50 {
		return error2.DescriptionTooLongError
	}
	_, parseErr := time.Parse("", request.AlarmStartDateTime)
	if parseErr != nil {
		return error2.InvalidAlarmDateTimeFormat
	}
	return nil
}

//TODO add logs
func (as alarmService) saveAlarm(ctx *gin.Context, createAlarmRequest request_model.CreateAlarmRequest) (string, *error2.ASError) {
	//TODO decide layout
	parsedTime, _ := time.Parse("", createAlarmRequest.AlarmStartDateTime)

	dbPrivateAlarmFlag := "F"
	if createAlarmRequest.Private {
		dbPrivateAlarmFlag = "T"
	}

	//TODO move this to UTIL else code becomes untestable.
	alarmID := uuid.New().String()
	transaction := as.transactionManager.NewTransaction()

	createAlarmDBError := as.alarmRepository.CreateAlarmMetadata(ctx, transaction, alarmID, createAlarmRequest.UserId, parsedTime, dbPrivateAlarmFlag, createAlarmRequest.Description)
	if createAlarmDBError != nil {
		transaction.Rollback()
		return "", error2.InternalServerError("error creating alarm")
	}

	var deviceAlarmSaveError error
	if createAlarmRequest.RepeatingDeviceAlarmIds.ContainsAtleastOneRepeatingAlarm() {
		deviceAlarmSaveError = as.saveRepeatingDeviceAlarmIds(ctx, transaction, createAlarmRequest.RepeatingDeviceAlarmIds, alarmID)
	} else {
		deviceAlarmSaveError = as.saveNonRepeatingDeviceAlarmId(ctx, transaction, *createAlarmRequest.NonRepeatingDeviceAlarmId, alarmID)
	}

	if deviceAlarmSaveError != nil {
		return "", error2.InternalServerError("could not save alarm.")
		transaction.Rollback()
	}
	commitError := transaction.Commit()
	if commitError != nil {
		return "", error2.InternalServerError("db commit failed.")
	}
	return alarmID, nil
}

func (as alarmService) saveRepeatingDeviceAlarmIds(ctx *gin.Context, transaction transaction_manager.Transaction, repeatingAlarmIds request_model.RepeatingDeviceAlarmIds, alarmId string) error {
	return as.alarmRepository.InsertRepeatingDeviceAlarmIDs(ctx, transaction, alarmId, repeatingAlarmIds.MapToDBModel())
}

func (as alarmService) saveNonRepeatingDeviceAlarmId(ctx *gin.Context, transaction transaction_manager.Transaction, deviceAlarmId int, alarmId string) error {
	return as.alarmRepository.InsertNonRepeatingDeviceAlarmID(ctx, transaction, alarmId, deviceAlarmId)
}

package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"social-alarm-service/constants"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/repository/transaction_manager"
	"social-alarm-service/request_model"
	"social-alarm-service/response_model"
	"time"
)

type AlarmService interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) (response_model.EligibleAlarmsResponse, *error2.ASError)
	CreateAlarm(ctx *gin.Context, request request_model.CreateAlarmRequest) (response_model.CreateAlarmResponse, *error2.ASError)
	UpdateStatus(ctx *gin.Context, alarmId string, userId string, status string) *error2.ASError
	GetAllAlarms(ctx *gin.Context, userId string) (response_model.GetAllAlarms, *error2.ASError)
	Delete(ctx *gin.Context, userId string, alarmId string) *error2.ASError
}

type alarmService struct {
	alarmRepository      repository.AlarmRepository
	userRepository       repository.UserRepository
	alarmMediaRepository repository.AlarmMediaRepository
	transactionManager   transaction_manager.TransactionManager
}

func NewAlarmService(alarmRepository repository.AlarmRepository, userRepository repository.UserRepository, alarmMediaRepository repository.AlarmMediaRepository,
	transactionManager transaction_manager.TransactionManager) AlarmService {
	return alarmService{alarmRepository: alarmRepository, userRepository: userRepository, alarmMediaRepository: alarmMediaRepository, transactionManager: transactionManager}
}

func (as alarmService) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) (response_model.EligibleAlarmsResponse, *error2.ASError) {
	fmt.Printf("fetching public non-expired alarms for user id %s \n", userId)

	userExists, repoErr := as.userRepository.UserExists(ctx, userId)
	if repoErr != nil {
		fmt.Printf("error when checking if user exists %v", repoErr)
		return response_model.EligibleAlarmsResponse{}, error2.InternalServerError("db error when checking if user exists")
	}
	if !userExists {
		fmt.Println("user does not exist in db")
		return response_model.EligibleAlarmsResponse{}, error2.InvalidUserIdError
	}

	publicNonExpiredRepeatingAlarms, publicNonExpiredNonRepeatingAlarms, err := as.alarmRepository.GetPublicNonExpiredAlarms(ctx, userId)
	if err != nil {
		fmt.Printf("error when fetching public non-expired alarms from db %s", err)
		return response_model.EligibleAlarmsResponse{}, error2.InternalServerError("db fetch error when getting public non expired alarms for given user id")
	}

	fmt.Printf("fetched %d public repeating alarms and %d non-repeating alarms \n", len(publicNonExpiredRepeatingAlarms), len(publicNonExpiredNonRepeatingAlarms))
	var eligibleAlarmResponse response_model.EligibleAlarmsResponse
	eligibleAlarmResponse.UserId = userId
	eligibleAlarmResponse.EligibleRepeatingAlarms = response_model.MapToEligibleRepeatingAlarms(publicNonExpiredRepeatingAlarms)
	eligibleAlarmResponse.EligibleNonRepeatingAlarms = response_model.MapToEligibleNonRepeatingAlarms(publicNonExpiredNonRepeatingAlarms)

	return eligibleAlarmResponse, nil
}

func (as alarmService) CreateAlarm(ctx *gin.Context, request request_model.CreateAlarmRequest) (response response_model.CreateAlarmResponse, asError *error2.ASError) {
	fmt.Println("Validating create alarm request")
	asError = as.validateCreateAlarmRequest(request)
	if asError != nil {
		fmt.Println("alarm creation request invalid. returning error")
		return
	}

	//TODO remove this once we have login mechanism in place. This api will be a protected post login API so we don't need to check if user exists.
	user, dbError := as.userRepository.GetUser(ctx, request.UserId)
	if dbError != nil && !errors.Is(dbError, sql.ErrNoRows) {
		if errors.Is(dbError, sql.ErrNoRows) {
			fmt.Printf("user with id %s does not exist", request.UserId)
			return response_model.CreateAlarmResponse{}, error2.InvalidUserIdError
		}
		fmt.Printf("error in db call to check if user exists in the db %v", dbError)
		asError = error2.InternalServerError("db fetch error")
		return
	}

	alarmId, asError := as.saveAlarm(ctx, request, user.CurrentTimezone)
	if asError != nil {
		return
	}

	fmt.Println("alarm saved successfully.")
	return response_model.CreateAlarmResponse{AlarmId: alarmId}, nil
}

func (as alarmService) UpdateStatus(ctx *gin.Context, alarmId string, userId string, status string) *error2.ASError {
	fmt.Printf("fetching alarm metadata for alarm id %s\n", alarmId)

	if status != constants.AlarmStatusON.String() && status != constants.AlarmStatusOFF.String() {
		return error2.InvalidAlarmStatus
	}

	alarms, dbErr := as.alarmRepository.GetAlarmMetadata(ctx, alarmId)
	if dbErr != nil {
		fmt.Printf("could not fetch alarm metadata for alarm id %s\n", alarmId)
		return error2.InternalServerError("db fetch failed")
	}

	if len(alarms) == 0 {
		fmt.Println("alarm id does not exist.")
		return error2.InvalidAlarmId
	}

	if alarms[0].UserID != userId {
		fmt.Println("alarm does not belong to the user in request.")
		return error2.OperationNotAllowed
	}

	dbErr = as.alarmRepository.UpdateAlarmStatus(ctx, alarmId, constants.GetAlarmStatus(status))
	if dbErr != nil {
		fmt.Printf("could not update status for alarm id %s . failed with error %v\n", alarmId, dbErr)
		return error2.InternalServerError("alarm update status failed")
	}
	fmt.Printf("alarm status updated successfully to %s\n", status)

	return nil
}

func (as alarmService) GetAllAlarms(ctx *gin.Context, userId string) (response_model.GetAllAlarms, *error2.ASError) {
	fmt.Printf("fetching all alarms for user id %s \n", userId)

	allAlarms := response_model.GetAllAlarms{}

	repeatingDbAlarms, err := as.alarmRepository.GetAllRepeatingAlarms(ctx, userId)
	if err != nil {
		fmt.Printf("error when fetching all repeating alarms for user id %s %v \n", userId, err)
		return response_model.GetAllAlarms{}, error2.InternalServerError("fetch failed")
	}
	fmt.Printf("user has %d repeating alarms \n", len(repeatingDbAlarms))

	ra := response_model.MapToRepeatingAlarms(repeatingDbAlarms)
	allAlarms.RepeatingAlarms = ra

	nonRepeatingDbAlarms, err := as.alarmRepository.GetAllNonRepeatingAlarms(ctx, userId)
	if err != nil {
		fmt.Printf("error when fetching all repeating alarms for user id %s %v \n", userId, err)
		return response_model.GetAllAlarms{}, error2.InternalServerError("non repeating db fetch failed")
	}
	fmt.Printf("user has %d non repeating alarms \n", len(nonRepeatingDbAlarms))

	nra := response_model.MapToNonRepeatingAlarms(nonRepeatingDbAlarms)
	allAlarms.NonRepeatingAlarms = nra

	return allAlarms, nil
}

func (as alarmService) Delete(ctx *gin.Context, userId string, alarmId string) *error2.ASError {
	//check if alarm is associated with user
	//delete from alarm media table
	//delete from repeating/non-repeating alarm table
	//delete from alarm table
	alarms, dbErr := as.alarmRepository.GetAlarmMetadata(ctx, alarmId)
	if dbErr != nil {
		fmt.Printf("could not fetch alarm metadata for alarm id %s\n", alarmId)
		return error2.InternalServerError("db fetch failed")
	}

	if len(alarms) == 0 {
		fmt.Println("alarm id does not exist.")
		return error2.InvalidAlarmId
	}

	if alarms[0].UserID != userId {
		fmt.Println("alarm does not belong to the user in request.")
		return error2.OperationNotAllowed
	}

	transaction := as.transactionManager.NewTransaction()

	dbErr = as.alarmMediaRepository.Delete(ctx, transaction, alarmId)
	if dbErr != nil {
		transaction.Rollback()
		return error2.InternalServerError("alarm media deletion failed")
	}

	dbErr = as.alarmRepository.DeleteRepeatingAlarm(ctx, transaction, alarmId)
	if dbErr != nil {
		transaction.Rollback()
		return error2.InternalServerError("repeating alarm  deletion failed")
	}

	dbErr = as.alarmRepository.DeleteNonRepeatingAlarm(ctx, transaction, alarmId)
	if dbErr != nil {
		transaction.Rollback()
		return error2.InternalServerError("non repeating alarm deletion failed")
	}

	dbErr = as.alarmRepository.Delete(ctx, transaction, alarmId)
	if dbErr != nil {
		transaction.Rollback()
		return error2.InternalServerError("alarm deletion failed")
	}

	err := transaction.Commit()
	if err != nil {
		return error2.InternalServerError("failed to delete alarm")
	}

	return nil
}

// TODO change layout. Request will include timezone as well. If alarm start datetime is an older date-time than current time , return error. Will add this check after all timezones are converted to UTC.
func (as alarmService) validateCreateAlarmRequest(request request_model.CreateAlarmRequest) *error2.ASError {
	if !request.RepeatingDeviceAlarmIds.ContainsAtleastOneRepeatingAlarm() && request.NonRepeatingDeviceAlarmId == nil {
		fmt.Println("request does not contain repeating or non repeating alarms. returning error")
		return error2.AlarmIdMissing
	}

	if request.RepeatingDeviceAlarmIds.ContainsAtleastOneRepeatingAlarm() && (request.NonRepeatingDeviceAlarmId != nil) {
		fmt.Println("request contains both repeating and non repeating alarms. returning error")
		return error2.InvalidAlarmTypeError
	}

	//TODO change layout. Request will include timezone as well.
	_, parseErr := time.Parse(constants.DateTimeLayout, request.AlarmStartDateTime)
	if parseErr != nil {
		fmt.Println("request has invalid date time format")
		return error2.InvalidAlarmDateTimeFormat
	}

	fmt.Println("create alarm request validated successfully")
	return nil
}

// TODO add logs
func (as alarmService) saveAlarm(ctx *gin.Context, createAlarmRequest request_model.CreateAlarmRequest, currentUserTimezone string) (string, *error2.ASError) {
	fmt.Println("saving alarm")
	//TODO decide layout
	parsedTime, _ := time.Parse("2006-01-02T15:04:05", createAlarmRequest.AlarmStartDateTime)

	location, _ := time.LoadLocation(currentUserTimezone)

	parsedTime = time.Date(
		parsedTime.Year(),
		parsedTime.Month(),
		parsedTime.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		parsedTime.Second(),
		parsedTime.Nanosecond(),
		location).UTC()

	alarmVisibility := constants.AlarmPrivateVisibility
	if !createAlarmRequest.Private {
		fmt.Println("alarm is private")
		alarmVisibility = constants.AlarmPublicVisibility
	}

	//TODO move this to UTIL else code becomes untestable.
	alarmID := uuid.New().String()
	transaction := as.transactionManager.NewTransaction()

	fmt.Printf("created uuid %s for saving alarm \n", alarmID)

	createAlarmDBError := as.alarmRepository.CreateAlarmMetadata(ctx, transaction, alarmID, createAlarmRequest.UserId, parsedTime, alarmVisibility, createAlarmRequest.Description)
	if createAlarmDBError != nil {
		transaction.Rollback()
		return "", error2.InternalServerError("error creating alarm")
	}

	fmt.Println("saved alarm metadata.")

	var deviceAlarmSaveError error
	if createAlarmRequest.RepeatingDeviceAlarmIds.ContainsAtleastOneRepeatingAlarm() {
		fmt.Println("alarm is repeating , saving repeating alarm")
		deviceAlarmSaveError = as.saveRepeatingDeviceAlarmIds(ctx, transaction, createAlarmRequest.RepeatingDeviceAlarmIds, alarmID)
	} else {
		fmt.Println("alarm is non repeating.")
		deviceAlarmSaveError = as.saveNonRepeatingDeviceAlarmId(ctx, transaction, *createAlarmRequest.NonRepeatingDeviceAlarmId, alarmID)
	}

	if deviceAlarmSaveError != nil {
		fmt.Printf("error when saving alarm , rolling back transaction %v \n", deviceAlarmSaveError)
		return "", error2.InternalServerError("could not save alarm.")
		transaction.Rollback()
	}

	commitError := transaction.Commit()
	if commitError != nil {
		fmt.Println("transaction commit error when saving transaction")
		return "", error2.InternalServerError("db commit failed.")
	}

	fmt.Println("alarm saved successfully in database")
	return alarmID, nil
}

func (as alarmService) saveRepeatingDeviceAlarmIds(ctx *gin.Context, transaction transaction_manager.Transaction, repeatingAlarmIds request_model.RepeatingDeviceAlarmIds, alarmId string) error {
	return as.alarmRepository.InsertRepeatingDeviceAlarmIDs(ctx, transaction, alarmId, repeatingAlarmIds.MapToDBModel())
}

func (as alarmService) saveNonRepeatingDeviceAlarmId(ctx *gin.Context, transaction transaction_manager.Transaction, deviceAlarmId int, alarmId string) error {
	return as.alarmRepository.InsertNonRepeatingDeviceAlarmID(ctx, transaction, alarmId, deviceAlarmId)
}

package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"mime/multipart"
	"os"
	"social-alarm-service/aws_util"
	"social-alarm-service/db_model"
	error2 "social-alarm-service/error"
	"social-alarm-service/repository"
	"social-alarm-service/repository/transaction_manager"
	"social-alarm-service/response_model"
)

type AlarmMediaService interface {
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError)
	UploadMedia(ctx *gin.Context, alarmId string, senderId string, file *multipart.File) *error2.ASError
}

type alarmMediaService struct {
	alarmRepo          repository.AlarmRepository
	alarmMediaRepo     repository.AlarmMediaRepository
	awsUtil            aws_util.AWSUtil
	transactionManager transaction_manager.TransactionManager
}

func NewAlarmMediaService(alarmRepo repository.AlarmRepository, alarmMediaRepo repository.AlarmMediaRepository, awsUtil aws_util.AWSUtil, transactionManager transaction_manager.TransactionManager) AlarmMediaService {
	return alarmMediaService{alarmRepo: alarmRepo, alarmMediaRepo: alarmMediaRepo, awsUtil: awsUtil, transactionManager: transactionManager}
}

func (as alarmMediaService) GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError) {
	alarmMedia, err := as.alarmMediaRepo.GetMediaForAlarm(ctx, alarmId)
	if err != nil {
		return []response_model.MediaForAlarm{}, error2.InternalServerError("db fetch error when getting all media associated with given alarm id")
	}
	return response_model.MapToMediaForAlarmResponseList(alarmMedia), nil
}

func (as alarmMediaService) UploadMedia(ctx *gin.Context, alarmId string, senderId string, file *multipart.File) (error *error2.ASError) {
	var alarm db_model.Alarms
	nonRepeatingAlarm, repoErr := as.alarmRepo.GetNonRepeatingAlarm(ctx, alarmId)
	if repoErr != nil {
		fmt.Println("error fetching non repeating alarm id")
		return error2.InternalServerError("error fetching non repeating alarm id")
	}
	if len(nonRepeatingAlarm) > 0 {
		alarm = nonRepeatingAlarm[0]
		if alarm.HasNonRepeatingAlarmExpired() || alarm.IsOff() || alarm.IsPrivate() {
			fmt.Println("media cannot be sent to this alarm as it has either expired , is private or is turned off")
			return error2.AlarmNotEligibleForMedia
		}
	}

	repeatingAlarm, repoErr := as.alarmRepo.GetRepeatingAlarm(ctx, alarmId)
	if repoErr != nil {
		return error2.InternalServerError("error fetching repeating alarm id")
	}
	if len(repeatingAlarm) > 0 {
		alarm = repeatingAlarm[0]
		if alarm.IsOff() || alarm.IsPrivate() {
			fmt.Println("media cannot be sent to this alarm as it is either private or turned off")
			return error2.AlarmNotEligibleForMedia
		}
	}

	if alarm.AlarmID == "" {
		fmt.Println("alarm id not found")
		return error2.InvalidAlarmId
	}

	fmt.Println("alarm is eligible to accept media. saving media file to aws")
	fileName, _ := uuid2.NewUUID()

	uploadError := as.awsUtil.UploadObject(ctx, file, os.Getenv("AWS_BUCKET_NAME"), fileName.String())
	if uploadError != nil {
		fmt.Printf("error when uploading resource to s3 %v \n", uploadError)
		return uploadError
	}

	defer func(fileName string) {
		if error != nil {
			as.awsUtil.DeleteObject(ctx, os.Getenv("AWS_BUCKET_NAME"), fileName)
		}
	}(fileName.String())

	transaction := as.transactionManager.NewTransaction()

	mediaId, _ := uuid2.NewUUID()
	resourceUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName.String())

	uploadMediaErr := as.alarmMediaRepo.UploadMedia(ctx, transaction, mediaId.String(), senderId, resourceUrl)
	if uploadMediaErr != nil {
		fmt.Printf("Error when creating media record %v \n", repoErr)
		transaction.Rollback()
		error = error2.InternalServerError("error when inserting media record")
		return
	}

	linkMediaErr := as.alarmMediaRepo.LinkMediaWithAlarm(ctx, transaction, alarmId, mediaId.String())
	if linkMediaErr != nil {
		fmt.Printf("error when linking media and alarm record %v \n", repoErr)
		transaction.Rollback()
		error = error2.InternalServerError("error when linking media and alarm record")
		return
	}

	commitErr := transaction.Commit()
	if commitErr != nil {
		fmt.Printf("error during commit %v \n", repoErr)
		transaction.Rollback()
		error = error2.InternalServerError("db commit error when saving media and linking media with alarm")
		return
	}

	fmt.Println("alarm media uploaded successfully")
	return nil
}

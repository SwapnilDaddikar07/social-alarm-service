package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"io/ioutil"
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
	UploadMedia(ctx *gin.Context, alarmId string, senderId string, fileName string) *error2.ASError
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]response_model.MediaForAlarm, *error2.ASError)
	CreateTmpFile(ctx *gin.Context, file multipart.File, extension string) (string, *error2.ASError)
	DeleteTmpFile(ctx *gin.Context, fileName string) *error2.ASError
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

//TODO check if this sender can send media to provided alarm i.e sender should be in contact of the receiver. Validation of sender id not needed as we will take it from token.
func (as alarmMediaService) UploadMedia(ctx *gin.Context, alarmId string, senderId string, fileName string) (error *error2.ASError) {
	fmt.Println("validating alarm id")
	error = as.validateAlarmId(ctx, alarmId)
	if error != nil {
		fmt.Println("error validating alarm id")
		return
	}

	fmt.Println("alarm is eligible to accept media. saving media file to aws")
	uploadError := as.awsUtil.UploadObject(ctx, "tmp/"+fileName, os.Getenv("AWS_BUCKET_NAME"), fileName)
	if uploadError != nil {
		fmt.Printf("error when uploading resource to s3 %v \n", uploadError)
		return uploadError
	}

	defer func(fileName string) {
		if error != nil {
			fmt.Println("removing saved object from s3 store as service threw error.")
			as.awsUtil.DeleteObject(ctx, os.Getenv("AWS_BUCKET_NAME"), fileName)
		}
	}(fileName)

	resourceUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)

	error = as.persistMediaMetadataAndLinkWithAlarm(ctx, alarmId, senderId, resourceUrl)
	if error != nil {
		fmt.Println("error when persisting alarm details to db")
		return
	}

	fmt.Println("alarm media uploaded successfully")
	return nil
}

func (as alarmMediaService) validateAlarmId(ctx *gin.Context, alarmId string) *error2.ASError {
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
	return nil
}

func (as alarmMediaService) CreateTmpFile(ctx *gin.Context, file multipart.File, extension string) (string, *error2.ASError) {
	b, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		fmt.Printf("could not read bytes from file %v ", readErr)
		return "", error2.InternalServerError("could not read bytes from request multipart file.")
	}

	tmpFileName, _ := uuid2.NewUUID()

	openFile, openErr := os.OpenFile("tmp/"+tmpFileName.String()+extension, os.O_CREATE|os.O_RDWR, 0660)
	if openErr != nil {
		fmt.Println("could not open file")
		return "", error2.InternalServerError("could not create tmp file.")
	}
	defer openFile.Close()

	_, writeErr := openFile.Write(b)
	if writeErr != nil {
		fmt.Println("could not write file", writeErr)
		return "", error2.InternalServerError("could not create tmp file.")
	}
	fmt.Printf("file %s successfully written to disk \n", tmpFileName.String())
	return tmpFileName.String() + extension, nil
}

func (as alarmMediaService) DeleteTmpFile(ctx *gin.Context, fileName string) *error2.ASError {
	deleteErr := os.Remove("tmp/" + fileName)
	if deleteErr != nil {
		fmt.Println("unable to delete tmp file ", deleteErr)
	} else {
		fmt.Printf("tmp file %s deleted successfully.\n", fileName)
	}
	return nil
}

func (as alarmMediaService) persistMediaMetadataAndLinkWithAlarm(ctx *gin.Context, alarmId string, senderId string, resourceUrl string) *error2.ASError {
	transaction := as.transactionManager.NewTransaction()

	mediaId, _ := uuid2.NewUUID()

	uploadMediaErr := as.alarmMediaRepo.UploadMedia(ctx, transaction, mediaId.String(), senderId, resourceUrl)
	if uploadMediaErr != nil {
		fmt.Printf("Error when creating media record %v \n", uploadMediaErr)
		transaction.Rollback()
		return error2.InternalServerError("error when inserting media record")
	}

	linkMediaErr := as.alarmMediaRepo.LinkMediaWithAlarm(ctx, transaction, alarmId, mediaId.String())
	if linkMediaErr != nil {
		fmt.Printf("error when linking media and alarm record %v \n", linkMediaErr)
		transaction.Rollback()
		return error2.InternalServerError("error when linking media and alarm record")
	}

	commitErr := transaction.Commit()
	if commitErr != nil {
		fmt.Printf("error during commit %v \n", commitErr)
		transaction.Rollback()
		return error2.InternalServerError("db commit error when saving media and linking media with alarm")
	}
	return nil
}

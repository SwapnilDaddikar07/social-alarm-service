package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"os"
	"social-alarm-service/db_model"
	error2 "social-alarm-service/error"
	"social-alarm-service/mocks"
	"testing"
	"time"
)

type AlarmMediaServiceTestSuite struct {
	suite.Suite
	context                *gin.Context
	mockCtrl               *gomock.Controller
	mockAlarmRepo          *mocks.MockAlarmRepository
	mockAlarmMediaRepo     *mocks.MockAlarmMediaRepository
	mockUserRepo           *mocks.MockUserRepository
	mockAWSUtil            *mocks.MockAWSUtil
	mockUtil               *mocks.MockUtils
	mockTransactionManager *mocks.MockTransactionManager
	mockTransaction        *mocks.MockTransaction
	alarmMediaService      AlarmMediaService
}

func TestAlarmMediaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AlarmMediaServiceTestSuite))
}

func (suite *AlarmMediaServiceTestSuite) SetupTest() {
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockAlarmRepo = mocks.NewMockAlarmRepository(suite.mockCtrl)
	suite.mockAlarmMediaRepo = mocks.NewMockAlarmMediaRepository(suite.mockCtrl)
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.mockCtrl)
	suite.mockAWSUtil = mocks.NewMockAWSUtil(suite.mockCtrl)
	suite.mockTransactionManager = mocks.NewMockTransactionManager(suite.mockCtrl)
	suite.mockTransaction = mocks.NewMockTransaction(suite.mockCtrl)
	suite.mockUtil = mocks.NewMockUtils(suite.mockCtrl)
	suite.alarmMediaService = NewAlarmMediaService(suite.mockAlarmRepo, suite.mockAlarmMediaRepo, suite.mockUserRepo, suite.mockAWSUtil, suite.mockUtil, suite.mockTransactionManager)
}

func (suite *AlarmMediaServiceTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_DBCallToCheckIfSenderExistsFails() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.InternalServerError("db error when checking if sender exists")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(false, errors.New("db error"))

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_SenderDoesNotExist() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.OperationNotAllowed

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(false, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_DBCallToCheckIfAlarmIsANonRepeatingAlarmFails() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.InternalServerError("error fetching non repeating alarm id")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, errors.New("db error"))

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_NonRepeatingAlarmIsPrivate() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.AlarmNotEligibleForMedia

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:    alarmId,
		UserID:     "owner-id",
		Visibility: "PRIVATE",
	}}, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_NonRepeatingAlarmIsExpired() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.AlarmNotEligibleForMedia

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:    alarmId,
		UserID:     "owner-id",
		Visibility: "PUBLIC",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, -2, 0),
		},
	}}, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_NonRepeatingAlarmIsTurnedOFF() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.AlarmNotEligibleForMedia

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:    alarmId,
		UserID:     "owner-id",
		Visibility: "PUBLIC",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, 1, 0),
		},
		Status: "OFF",
	}}, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_DBCallToCheckIfAlarmIsRepeatingFails() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.InternalServerError("error fetching repeating alarm id")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, nil)
	suite.mockAlarmRepo.EXPECT().GetRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, errors.New("db error"))

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_RepeatingAlarmIsPrivate() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.AlarmNotEligibleForMedia

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, nil)
	suite.mockAlarmRepo.EXPECT().GetRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:                   alarmId,
		UserID:                    "owner-id",
		Visibility:                "PRIVATE",
		NonRepeatingDeviceAlarmId: 1,
		MonDeviceAlarmId:          3,
		TueDeviceAlarmId:          0,
		WedDeviceAlarmId:          0,
		ThuDeviceAlarmId:          0,
		FriDeviceAlarmId:          0,
		SatDeviceAlarmId:          0,
		SunDeviceAlarmId:          2,
	}}, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_RepeatingAlarmIsOFF() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.AlarmNotEligibleForMedia

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, nil)
	suite.mockAlarmRepo.EXPECT().GetRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:                   alarmId,
		UserID:                    "owner-id",
		Visibility:                "PUBLIC",
		Status:                    "OFF",
		NonRepeatingDeviceAlarmId: 1,
		MonDeviceAlarmId:          3,
		TueDeviceAlarmId:          0,
		WedDeviceAlarmId:          0,
		ThuDeviceAlarmId:          0,
		FriDeviceAlarmId:          0,
		SatDeviceAlarmId:          0,
		SunDeviceAlarmId:          2,
	}}, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_AlarmIdIsNotPresentInDB() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	expectedError := error2.InvalidAlarmId

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, nil)
	suite.mockAlarmRepo.EXPECT().GetRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{}, nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_FileUploadToAWSFails() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	bucketName := "test-bucket-name"
	os.Setenv("AWS_BUCKET_NAME", bucketName)
	expectedError := error2.InternalServerError("aws upload failed")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:     alarmId,
		UserID:      "owner-id",
		Visibility:  "PUBLIC",
		Description: "some-description",
		Status:      "ON",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, 1, 0),
		},
		CreatedAt:                 sql.NullTime{},
		NonRepeatingDeviceAlarmId: 1,
	}}, nil)

	suite.mockAWSUtil.EXPECT().UploadObject(suite.context, "tmp/"+fileName, bucketName, fileName).Return(error2.InternalServerError("aws upload failed"))

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_UploadToMediaTableInDBFails_AndShouldDeleteObjectFromAWS() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	bucketName := "test-bucket-name"
	awsRegion := "asia"
	mediaId := "media-id"
	os.Setenv("AWS_BUCKET_NAME", bucketName)
	os.Setenv("AWS_REGION", awsRegion)
	resourceUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)

	expectedError := error2.InternalServerError("error when inserting media record")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:     alarmId,
		UserID:      "owner-id",
		Visibility:  "PUBLIC",
		Description: "some-description",
		Status:      "ON",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, 1, 0),
		},
		CreatedAt:                 sql.NullTime{},
		NonRepeatingDeviceAlarmId: 1,
	}}, nil)

	suite.mockAWSUtil.EXPECT().UploadObject(suite.context, "tmp/"+fileName, bucketName, fileName).Return(nil)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction)
	suite.mockUtil.EXPECT().GenerateUUID().Return(mediaId)
	suite.mockAlarmMediaRepo.EXPECT().UploadMedia(suite.context, suite.mockTransaction, mediaId, senderId, resourceUrl).Return(errors.New("db error"))
	suite.mockTransaction.EXPECT().Rollback().Return(nil)
	suite.mockAWSUtil.EXPECT().DeleteObject(suite.context, bucketName, fileName).Return(nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_LinkingMediaWithAlarmIdInDBFails_AndShouldDeleteObjectFromAWS() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	bucketName := "test-bucket-name"
	awsRegion := "asia"
	mediaId := "media-id"
	os.Setenv("AWS_BUCKET_NAME", bucketName)
	os.Setenv("AWS_REGION", awsRegion)
	resourceUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)

	expectedError := error2.InternalServerError("error when linking media and alarm record")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:     alarmId,
		UserID:      "owner-id",
		Visibility:  "PUBLIC",
		Description: "some-description",
		Status:      "ON",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, 1, 0),
		},
		CreatedAt:                 sql.NullTime{},
		NonRepeatingDeviceAlarmId: 1,
	}}, nil)

	suite.mockAWSUtil.EXPECT().UploadObject(suite.context, "tmp/"+fileName, bucketName, fileName).Return(nil)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction)
	suite.mockUtil.EXPECT().GenerateUUID().Return(mediaId)
	suite.mockAlarmMediaRepo.EXPECT().UploadMedia(suite.context, suite.mockTransaction, mediaId, senderId, resourceUrl).Return(nil)
	suite.mockAlarmMediaRepo.EXPECT().LinkMediaWithAlarm(suite.context, suite.mockTransaction, alarmId, mediaId).Return(errors.New("db error"))
	suite.mockTransaction.EXPECT().Rollback().Return(nil)
	suite.mockAWSUtil.EXPECT().DeleteObject(suite.context, bucketName, fileName).Return(nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_ShouldThrowError_When_TransactionCommitFails_AndDeleteObjectUploadedToAWS() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	bucketName := "test-bucket-name"
	awsRegion := "asia"
	mediaId := "media-id"
	os.Setenv("AWS_BUCKET_NAME", bucketName)
	os.Setenv("AWS_REGION", awsRegion)
	resourceUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)

	expectedError := error2.InternalServerError("db commit error when saving media and linking media with alarm")

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:     alarmId,
		UserID:      "owner-id",
		Visibility:  "PUBLIC",
		Description: "some-description",
		Status:      "ON",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, 1, 0),
		},
		CreatedAt:                 sql.NullTime{},
		NonRepeatingDeviceAlarmId: 1,
	}}, nil)

	suite.mockAWSUtil.EXPECT().UploadObject(suite.context, "tmp/"+fileName, bucketName, fileName).Return(nil)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction)
	suite.mockUtil.EXPECT().GenerateUUID().Return(mediaId)
	suite.mockAlarmMediaRepo.EXPECT().UploadMedia(suite.context, suite.mockTransaction, mediaId, senderId, resourceUrl).Return(nil)
	suite.mockAlarmMediaRepo.EXPECT().LinkMediaWithAlarm(suite.context, suite.mockTransaction, alarmId, mediaId).Return(nil)
	suite.mockTransaction.EXPECT().Commit().Return(errors.New("commit error"))
	suite.mockTransaction.EXPECT().Rollback().Return(nil)
	suite.mockAWSUtil.EXPECT().DeleteObject(suite.context, bucketName, fileName).Return(nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Equal(expectedError, actualErr)
}

func (suite *AlarmMediaServiceTestSuite) TestAlarmMediaService_UploadMedia_NotThrowError_When_MediaIsSuccessfullyStoredInDBAndUploadedToAWS() {
	alarmId := "some-alarm-id"
	senderId := "some-sender-id"
	fileName := "file-name"
	bucketName := "test-bucket-name"
	awsRegion := "asia"
	mediaId := "media-id"
	os.Setenv("AWS_BUCKET_NAME", bucketName)
	os.Setenv("AWS_REGION", awsRegion)
	resourceUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)

	suite.mockUserRepo.EXPECT().UserExists(suite.context, senderId).Return(true, nil)
	suite.mockAlarmRepo.EXPECT().GetNonRepeatingAlarm(suite.context, alarmId).Return([]db_model.Alarms{{
		AlarmID:     alarmId,
		UserID:      "owner-id",
		Visibility:  "PUBLIC",
		Description: "some-description",
		Status:      "ON",
		AlarmStartDateTime: sql.NullTime{
			Time: time.Now().AddDate(0, 1, 0),
		},
		CreatedAt:                 sql.NullTime{},
		NonRepeatingDeviceAlarmId: 1,
	}}, nil)

	suite.mockAWSUtil.EXPECT().UploadObject(suite.context, "tmp/"+fileName, bucketName, fileName).Return(nil)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction)
	suite.mockUtil.EXPECT().GenerateUUID().Return(mediaId)
	suite.mockAlarmMediaRepo.EXPECT().UploadMedia(suite.context, suite.mockTransaction, mediaId, senderId, resourceUrl).Return(nil)
	suite.mockAlarmMediaRepo.EXPECT().LinkMediaWithAlarm(suite.context, suite.mockTransaction, alarmId, mediaId).Return(nil)
	suite.mockTransaction.EXPECT().Commit().Return(nil)

	actualErr := suite.alarmMediaService.UploadMedia(suite.context, alarmId, senderId, fileName)

	suite.Nil(nil, actualErr)
}

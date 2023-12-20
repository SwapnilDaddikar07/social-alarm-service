package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/db_model"
	"social-alarm-service/repository/transaction_manager"
)

type AlarmMediaRepository interface {
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]db_model.GetMediaForAlarm, error)
	UploadMedia(ctx *gin.Context, transaction transaction_manager.Transaction, mediaId string, senderId string, mediaURL string) error
	LinkMediaWithAlarm(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID, mediaID string) error
	Delete(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string) error
}

type alarmMediaRepository struct {
	db *sqlx.DB
}

func NewAlarmMediaRepository(db *sqlx.DB) AlarmMediaRepository {
	return alarmMediaRepository{db: db}
}

func (ar alarmMediaRepository) GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]db_model.GetMediaForAlarm, error) {
	mediaForAlarms := make([]db_model.GetMediaForAlarm, 0)
	query := "select u.display_name , m.resource_url from media m  inner join alarm_media am on m.media_id = am.media_id  inner join users u on u.user_id = m.sender_id  where am.alarm_id = ? order by m.created_at desc"

	dbFetchError := ar.db.Select(&mediaForAlarms, query, alarmId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting all media for alarm id", dbFetchError)
		return mediaForAlarms, dbFetchError
	}
	return mediaForAlarms, dbFetchError
}

func (ar alarmMediaRepository) UploadMedia(ctx *gin.Context, transaction transaction_manager.Transaction, mediaId string, senderId string, resourceURL string) error {
	query := "insert into media (media_id , sender_id , resource_url , created_at) values (?,?,?, CURRENT_TIME)"

	_, dbError := transaction.Exec(query, mediaId, senderId, resourceURL)
	return dbError
}

func (ar alarmMediaRepository) LinkMediaWithAlarm(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID, mediaID string) error {
	query := "insert into alarm_media (alarm_id , media_id ) values (?,?)"

	_, dbError := transaction.Exec(query, alarmID, mediaID)
	return dbError
}

func (ar alarmMediaRepository) Delete(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string) error {
	query := "delete from alarm_media where alarm_id = ?"

	_, dbError := transaction.Exec(query, alarmId)
	return dbError
}

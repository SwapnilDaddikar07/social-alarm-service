package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/db_model"
)

type AlarmMediaRepository interface {
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]db_model.GetMediaForAlarm, error)
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

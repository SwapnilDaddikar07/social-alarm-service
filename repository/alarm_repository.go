package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/db_model"
)

type AlarmRepository interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredAlarms, error)
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]db_model.GetMediaForAlarm, error)
}

type alarmRepository struct {
	db *sqlx.DB
}

func NewAlarmRepository(db *sqlx.DB) AlarmRepository {
	return alarmRepository{db: db}
}

func (ar alarmRepository) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredAlarms, error) {
	publicNonExpiredAlarms := make([]db_model.PublicNonExpiredAlarms, 0)
	query := "select a.alarm_id , a.alarm_start_datetime , a.description , ash.mon , ash.tue ,ash.wed , ash.thu, ash.fri ,ash.sat, ash.sun FROM alarms a inner join alarm_schedules ash on a.alarm_id = ash.alarm_id where a.user_id = ? and a.status = 'ON' and a.visibility = 'P' AND (a.type = 'R' OR (a.alarm_start_datetime > CURRENT_TIMESTAMP))"

	dbFetchError := ar.db.Select(&publicNonExpiredAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error while getting all public non expired alarms for user id", dbFetchError)
		return publicNonExpiredAlarms, dbFetchError
	}
	return publicNonExpiredAlarms, dbFetchError
}

func (ar alarmRepository) GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]db_model.GetMediaForAlarm, error) {
	mediaForAlarms := make([]db_model.GetMediaForAlarm, 0)
	query := "select u.display_name , m.resource_url from alarm_media am  inner join media m on am.media_id = m.media_id  inner join users u on u.user_id = m.sender_id  where am.alarm_id = ? order by m.created_at desc"

	dbFetchError := ar.db.Select(&mediaForAlarms, query, alarmId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting all media for alarm id", dbFetchError)
		return mediaForAlarms, dbFetchError
	}
	return mediaForAlarms, dbFetchError
}

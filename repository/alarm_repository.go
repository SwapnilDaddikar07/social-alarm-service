package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/db_model"
)

type AlarmRepository interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredAlarms, error)
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
	return publicNonExpiredAlarms, nil
}

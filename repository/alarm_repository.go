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

/*
	Final query
	select alarms.alarm_id , alarms.alarm_start_date_time , alarms.description , alarm_schedules.Mon , alarm_schedules.Tue ,alarm_schedules.Wed , alarm_schedules.Thu, alarm_schedules.Fri ,alarm_schedules.Sat, alarm_schedules.Sun FROM alarms inner join alarms_schedules on alarms.alarm_id = alarm_schedules.alarm_id where user_id = ? and status = 'ON' and visibility = 'P' AND (type == 'R' OR (alarm_start_date_time > CURRENT_TIMESTAMP))
*/
//GetAllAlarms TODO Change the response model to a model specific to DB as we want to parse schedules. Replace with Final query mentioned above.
func (ar alarmRepository) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredAlarms, error) {
	publicNonExpiredAlarms := make([]db_model.PublicNonExpiredAlarms, 0)
	query := "select alarm_id , alarm_start_datetime , description where user_id = ? and status = 'ON' and visibility = 'P'"

	dbFetchError := ar.db.Select(&publicNonExpiredAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db error", dbFetchError)
		return publicNonExpiredAlarms, dbFetchError
	}
	return publicNonExpiredAlarms, nil
}

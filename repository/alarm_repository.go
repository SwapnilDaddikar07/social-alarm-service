package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/model"
)

type AlarmRepository interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error)
}

type alarmRepository struct {
	db *sqlx.DB
}

func NewAlarmRepository(db *sqlx.DB) AlarmRepository {
	return alarmRepository{db: db}
}

/*
	Final query
	select alarms.alarm_id , alarms.alarm_start_date_time , alarms.description , alarm_schedules.M , alarm_schedules.Tue ,alarm_schedules.W , alarm_schedules.Thu, alarm_schedules.F,alarm_schedules.Sat, alarm_schedules.Sun FROM alarms inner join alarms_schedules on alarms.alarm_id = alarm_schedules.alarm_id where user_id =:userId and status = 'ON' and visibility = 'P' AND (type == 'R' OR (alarm_start_date_time > CURRENT_TIMESTAMP))
*/
//GetAllAlarms TODO Change the response model to a model specific to DB as we want to parse schedules. Replace with Final query mentioned above.
func (ar alarmRepository) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error) {
	response := make([]model.AlarmResponse, 0)
	query := "select alarms.alarm_id , alarms.alarm_start_date_time , alarms.description where user_id =:userId and status = 'ON' and visibility = 'P' AND (type == 'R' OR (alarm_start_date_time > CURRENT_TIMESTAMP))"

	dbErr := ar.db.Select(&response, query, map[string]interface{}{
		"userId": userId,
	})
	if dbErr != nil {
		fmt.Println("db error", dbErr)
		return response, dbErr
	}
	return response, nil
}

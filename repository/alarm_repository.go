package repository

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"social-alarm-service/model"
)

type AlarmRepository interface {
	GetAllAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error)
}

type alarmRepository struct {
	db sql.DB
}

func NewAlarmRepository(db sql.DB) AlarmRepository {
	return alarmRepository{db: db}
}

func (ar alarmRepository) GetAllAlarms(ctx *gin.Context, userId string) ([]model.AlarmResponse, error) {
	result := make([]model.AlarmResponse, 0)

	rows, _ := ar.db.Query("select alarms.alarm_id , alarms.alarm_start_date_time , alarms.description , alarm_schedules.M , alarm_schedules.Tue ,alarm_schedules.W , alarm_schedules.Thu, alarm_schedules.F,alarm_schedules.Sat, alarm_schedules.Sun  from alarms a where user_id = userId and status = 'ON' and visibility = 'P' AND (type == 'R' OR (alarm_start_date_time > currentTime)) inner join alarms_schedules on alarms.alarm_id = alarm_schedules.alarm_id")
	for rows.Next() {
		alarmResponse := &model.AlarmResponse{}
		err := rows.Scan(alarmResponse)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *alarmResponse)
	}
	return result, nil
}

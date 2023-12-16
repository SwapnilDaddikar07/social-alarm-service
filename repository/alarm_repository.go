package repository

import (
	"errors"
	"fmt"
	"social-alarm-service/constants"
	"social-alarm-service/db_model"
	"social-alarm-service/repository/transaction_manager"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AlarmRepository interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, []db_model.Alarms, error)
	GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error)
	GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error)
	CreateAlarmMetadata(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string, userId string, alarmStartDateTime time.Time, alarmType constants.AlarmVisibility, description string) error
	InsertNonRepeatingDeviceAlarmID(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, deviceAlarmID int) error
	InsertRepeatingDeviceAlarmIDs(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, repeatingIDs db_model.RepeatingAlarmIDs) error
	GetRepeatingAlarm(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error)
	GetNonRepeatingAlarm(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error)
	UpdateAlarmStatus(ctx *gin.Context, alarmId string, status constants.AlarmStatus) error
	GetAlarmMetadata(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error)
	GetAllRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error)
	GetAllNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error)
}

type alarmRepository struct {
	db *sqlx.DB
}

func NewAlarmRepository(db *sqlx.DB) AlarmRepository {
	return alarmRepository{db: db}
}

func (ar alarmRepository) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) (repeatingAlarms []db_model.Alarms, nonRepeatingAlarms []db_model.Alarms, err error) {
	repeatingAlarms, err = ar.GetPublicNonExpiredRepeatingAlarms(ctx, userId)
	if err != nil {
		return
	}
	nonRepeatingAlarms, err = ar.GetPublicNonExpiredNonRepeatingAlarms(ctx, userId)
	if err != nil {
		return
	}
	fmt.Printf("successfully fetched repeating and non repeating alarms for user id %s \n", userId)
	return repeatingAlarms, nonRepeatingAlarms, nil
}

func (ar alarmRepository) GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	repeatingAlarms := make([]db_model.Alarms, 0)

	query := "select a.alarm_id , a.alarm_start_datetime , a.description , rda.mon_device_alarm_id , rda.tue_device_alarm_id , rda.wed_device_alarm_id , rda.thu_device_alarm_id , rda.fri_device_alarm_id , rda.sat_device_alarm_id , rda.sun_device_alarm_id " +
		"from alarms a inner join repeating_device_alarm_id rda " +
		"on a.alarm_id = rda.alarm_id where a.visibility = 'PUBLIC' and a.status='ON' and a.user_id= ?"

	dbFetchError := ar.db.Select(&repeatingAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting public non expired repeating alarms", dbFetchError)
		return repeatingAlarms, dbFetchError
	}

	return repeatingAlarms, dbFetchError
}

func (ar alarmRepository) GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	nonRepeatingAlarms := make([]db_model.Alarms, 0)

	query := "select a.alarm_id , a.alarm_start_datetime , a.description " +
		"from alarms a inner join non_repeating_device_alarm_id nrda " +
		"on a.alarm_id = nrda.alarm_id where a.user_id= ? and a.visibility = 'PUBLIC' and a.status='ON' and a.alarm_start_datetime > CURRENT_TIME"

	dbFetchError := ar.db.Select(&nonRepeatingAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting public non expired non repeating alarms", dbFetchError)
		return nonRepeatingAlarms, dbFetchError
	}

	return nonRepeatingAlarms, dbFetchError
}

func (ar alarmRepository) CreateAlarmMetadata(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string, userId string, alarmStartDateTime time.Time, visibility constants.AlarmVisibility, description string) error {
	query := "INSERT INTO alarms (alarm_id, user_id , alarm_start_datetime , created_at , visibility , description, status) " +
		"VALUES " +
		"(?,?,?,CURRENT_TIME,?,?,?)"

	_, dbError := transaction.Exec(query, alarmId, userId, alarmStartDateTime, visibility, description, "ON")
	return dbError
}

func (ar alarmRepository) InsertNonRepeatingDeviceAlarmID(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, deviceAlarmID int) error {
	query := "INSERT INTO non_repeating_device_alarm_id (alarm_id , device_alarm_id) " +
		"VALUES " +
		"(?,?)"

	_, dbError := transaction.Exec(query, alarmID, deviceAlarmID)
	return dbError
}

func (ar alarmRepository) InsertRepeatingDeviceAlarmIDs(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, repeatingIDs db_model.RepeatingAlarmIDs) error {

	query := "INSERT INTO repeating_device_alarm_id (alarm_id , mon_device_alarm_id , tue_device_alarm_id , wed_device_alarm_id , thu_device_alarm_id , fri_device_alarm_id , sat_device_alarm_id , sun_device_alarm_id) " +
		"VALUES" +
		"?,?,?,?,?,?,?,?"

	_, dbError := transaction.Exec(query, alarmID, repeatingIDs.Mon, repeatingIDs.Tue, repeatingIDs.Wed, repeatingIDs.Thu, repeatingIDs.Fri, repeatingIDs.Sat, repeatingIDs.Sun)
	return dbError
}

func (ar alarmRepository) GetRepeatingAlarm(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error) {
	query := "select a.alarm_id, a.user_id, a.visibility, a.description, a.status, a.created_at, a.alarm_start_datetime, " +
		"rda.mon_device_alarm_id, rda.tue_device_alarm_id, rda.wed_device_alarm_id, rda.thu_device_alarm_id, rda.fri_device_alarm_id, rda.sat_device_alarm_id, rda.sun_device_alarm_id " +
		"from alarms a " +
		"join repeating_device_alarm_id rda " +
		"on a.alarm_id = rda.alarm_id where a.alarm_id = ?"

	var alarms []db_model.Alarms
	dbErr := ar.db.Select(&alarms, query, alarmId)
	if dbErr != nil {
		return []db_model.Alarms{}, dbErr
	}
	return alarms, nil
}

func (ar alarmRepository) GetNonRepeatingAlarm(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error) {
	query := "select a.alarm_id , a.user_id , a.visibility , a.description , a.status , a.created_at , a.alarm_start_datetime from alarms a " +
		"join non_repeating_device_alarm_id nrda " +
		"on a.alarm_id = nrda.alarm_id where a.alarm_id = ?"

	var alarms []db_model.Alarms
	dbErr := ar.db.Select(&alarms, query, alarmId)
	if dbErr != nil {
		return []db_model.Alarms{}, dbErr
	}
	return alarms, nil
}

func (ar alarmRepository) GetAlarmMetadata(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error) {
	query := "select alarm_id, user_id, visibility, description, status, created_at, alarm_start_datetime from alarms where alarm_id = ?"

	var alarms []db_model.Alarms
	dbErr := ar.db.Select(&alarms, query, alarmId)
	if dbErr != nil {
		fmt.Println("db error", dbErr)
		return []db_model.Alarms{}, dbErr
	}
	return alarms, nil
}

func (ar alarmRepository) UpdateAlarmStatus(ctx *gin.Context, alarmId string, status constants.AlarmStatus) error {
	query := "update alarms set status=? where alarm_id=?"

	_, err := ar.db.Exec(query, status, alarmId)
	if err != nil {
		fmt.Printf("db operation to update status failed for alarm id %s \n", alarmId)
		return errors.New("update status failed")
	}
	return nil
}

func (ar alarmRepository) GetAllRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	query := "select a.alarm_id, a.user_id, a.visibility, a.description, a.status, a.created_at, a.alarm_start_datetime, " +
		"rda.mon_device_alarm_id, rda.tue_device_alarm_id, rda.wed_device_alarm_id, rda.thu_device_alarm_id , rda.fri_device_alarm_id, " +
		"rda.sat_device_alarm_id, rda.sun_device_alarm_id, count(am.alarm_id) as media_count" +
		"from alarms a " +
		"inner join repeating_device_alarm_id rda on a.alarm_id = rda.alarm_id " +
		"left join alarm_media am on am.alarm_id = a.alarm_id" +
		"where a.user_id = ? " +
		"group by a.alarm_id"

	var alarms []db_model.Alarms
	dbErr := ar.db.Select(&alarms, query, userId)
	if dbErr != nil {
		fmt.Println("db error", dbErr)
		return []db_model.Alarms{}, dbErr
	}
	return alarms, nil
}

func (ar alarmRepository) GetAllNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	query := "select a.alarm_id, a.user_id, a.visibility, a.description, a.status, a.created_at, a.alarm_start_datetime, " +
		"nrda.device_alarm_id, count(am.alarm_id) as media_count" +
		"from alarms a " +
		"inner join non_repeating_device_alarm_id nrda on a.alarm_id = nrda.alarm_id " +
		"left join alarm_media am on a.alarm_id = am.alarm_id" +
		"where a.user_id = ? " +
		"group by a.alarm_id"

	var alarms []db_model.Alarms
	dbErr := ar.db.Select(&alarms, query, userId)
	if dbErr != nil {
		fmt.Println("db error", dbErr)
		return []db_model.Alarms{}, dbErr
	}
	return alarms, nil
}

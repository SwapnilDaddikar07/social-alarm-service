package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/constants"
	"social-alarm-service/db_model"
	"social-alarm-service/repository/transaction_manager"
	"time"
)

type AlarmRepository interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, []db_model.PublicNonExpiredNonRepeatingAlarms, error)
	GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, error)
	GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredNonRepeatingAlarms, error)
	UserExists(ctx *gin.Context, userId string) (bool, error)
	CreateAlarmMetadata(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string, userId string, alarmStartDateTime time.Time, alarmType constants.AlarmVisibility, description string) error
	InsertNonRepeatingDeviceAlarmID(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, deviceAlarmID int) error
	InsertRepeatingDeviceAlarmIDs(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, repeatingIDs db_model.RepeatingAlarmIDs) error
}

type alarmRepository struct {
	db *sqlx.DB
}

func NewAlarmRepository(db *sqlx.DB) AlarmRepository {
	return alarmRepository{db: db}
}

func (ar alarmRepository) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, []db_model.PublicNonExpiredNonRepeatingAlarms, error) {
	repeatingAlarms, dbError := ar.GetPublicNonExpiredRepeatingAlarms(ctx, userId)
	if dbError != nil {
		return []db_model.PublicNonExpiredRepeatingAlarms{}, []db_model.PublicNonExpiredNonRepeatingAlarms{}, dbError
	}

	nonRepeatingAlarms, dbError := ar.GetPublicNonExpiredNonRepeatingAlarms(ctx, userId)
	if dbError != nil {
		return []db_model.PublicNonExpiredRepeatingAlarms{}, []db_model.PublicNonExpiredNonRepeatingAlarms{}, dbError
	}
	fmt.Printf("successfully fetched repeating and non repeating alarms for user id %s \n", userId)
	return repeatingAlarms, nonRepeatingAlarms, nil
}

func (ar alarmRepository) GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, error) {
	mediaForAlarms := make([]db_model.PublicNonExpiredRepeatingAlarms, 0)

	query := "select a.alarm_id , a.alarm_start_datetime , a.description , rda.mon_device_alarm_id , rda.tue_device_alarm_id , rda.wed_device_alarm_id , rda.thu_device_alarm_id , rda.fri_device_alarm_id , rda.sat_device_alarm_id , rda.sun_device_alarm_id " +
		"from alarms a inner join repeating_device_alarm_id rda " +
		"on a.alarm_id = rda.alarm_id where a.visibility = 'PUBLIC' and a.status='ON' and a.user_id= ?"

	dbFetchError := ar.db.Select(&mediaForAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting public non expired repeating alarms", dbFetchError)
		return mediaForAlarms, dbFetchError
	}

	return mediaForAlarms, dbFetchError
}

func (ar alarmRepository) GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredNonRepeatingAlarms, error) {
	mediaForAlarms := make([]db_model.PublicNonExpiredNonRepeatingAlarms, 0)

	query := "select a.alarm_id , a.alarm_start_datetime , a.description " +
		"from alarms a inner join non_repeating_device_alarm_id nrda " +
		"on a.alarm_id = nrda.alarm_id where a.user_id= ? and a.visibility = 'PUBLIC' and a.status='ON' and a.alarm_start_datetime > CURRENT_TIME"

	dbFetchError := ar.db.Select(&mediaForAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting public non expired non repeating alarms", dbFetchError)
		return mediaForAlarms, dbFetchError
	}
	return mediaForAlarms, dbFetchError
}

func (ar alarmRepository) UserExists(ctx *gin.Context, userId string) (bool, error) {
	query := "SELECT EXISTS(SELECT user_id from users WHERE user_id= ?)"
	var rows []int

	dbFetchError := ar.db.Select(&rows, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when checking if user id exists in the db", dbFetchError)
		return false, dbFetchError
	}
	return len(rows) == 1, nil
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

package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/db_model"
	"social-alarm-service/repository/transaction_manager"
	"time"
)

type AlarmRepository interface {
	GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, []db_model.PublicNonExpiredNonRepeatingAlarms, error)
	GetMediaForAlarm(ctx *gin.Context, alarmId string) ([]db_model.GetMediaForAlarm, error)
	GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, error)
	GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredNonRepeatingAlarms, error)
	UserExists(ctx *gin.Context, userId string) (bool, error)
	CreateAlarmMetadata(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string, userId string, alarmStartDateTime time.Time, isPrivate string, description string) error
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

func (ar alarmRepository) GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredRepeatingAlarms, error) {
	mediaForAlarms := make([]db_model.PublicNonExpiredRepeatingAlarms, 0)

	query := "select a.alarm_id , a.alarm_start_datetime , a.description , rsa.mon_system_alarm_id , rsa.tue_system_alarm_id , rsa.wed_system_alarm_id , rsa.thu_system_alarm_id , rsa.fri_system_alarm_id , rsa.sat_system_alarm_id , rsa.sun_system_alarm_id " +
		"from alarm a inner join repeating_system_alarm_id rsa " +
		"on a.alarm_id = rsa.alarm_id where a.user_id= ?"

	dbFetchError := ar.db.Select(&mediaForAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting public non expired repeating alarms", dbFetchError)
		return mediaForAlarms, dbFetchError
	}

	return mediaForAlarms, dbFetchError
}

func (ar alarmRepository) GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.PublicNonExpiredNonRepeatingAlarms, error) {
	mediaForAlarms := make([]db_model.PublicNonExpiredNonRepeatingAlarms, 0)

	query := "select a.alarm_id , a.alarm_start_datetime , a.description , ssa.system_alarm_id " +
		"from alarm a inner join non_repeating_system_alarm_id ssa " +
		"on a.alarm_id = ssa.alarm_id where a.user_id= ? and a.alarm_start_datetime > CURRENT_TIME"

	dbFetchError := ar.db.Select(&mediaForAlarms, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when getting public non expired non repeating alarms", dbFetchError)
		return mediaForAlarms, dbFetchError
	}
	return mediaForAlarms, dbFetchError
}

func (ar alarmRepository) UserExists(ctx *gin.Context, userId string) (bool, error) {
	query := "SELECT EXISTS(SELECT user_id from users WHERE user_id= ?)"
	var rows *int

	dbFetchError := ar.db.Select(rows, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when checking if user id exists in the db", dbFetchError)
		return false, dbFetchError
	}
	return *rows == 1, nil
}

// CreateAlarmMetadata TODO change column name from system alarm id to device alarm id
func (ar alarmRepository) CreateAlarmMetadata(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId string, userId string, alarmStartDateTime time.Time, isPrivate string, description string) error {
	query := "INSERT INTO alarms (alarm_id, user_id , alarm_start_date_time , visibility , description, status) " +
		"VALUES " +
		"(?,?,?,?,?,?)"

	_, dbError := transaction.Exec(query, alarmId, userId, alarmStartDateTime, isPrivate, description, "ON")
	return dbError
}

// InsertNonRepeatingDeviceAlarmID TODO change column name from system alarm id to device alarm id
func (ar alarmRepository) InsertNonRepeatingDeviceAlarmID(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, deviceAlarmID int) error {
	query := "INSERT INTO non_repeating_system_alarm_id (alarm_id , repeating_system_alarm_id) " +
		"VALUES " +
		"(?,?)"

	_, dbError := transaction.Exec(query, alarmID, deviceAlarmID)
	return dbError
}

// InsertRepeatingDeviceAlarmIDs TODO change column name from system alarm id to device alarm id
func (ar alarmRepository) InsertRepeatingDeviceAlarmIDs(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, repeatingIDs db_model.RepeatingAlarmIDs) error {

	query := "INSERT INTO repeating_system_alarm_id (alarm_id , mon , tue , wed , thu , fri , sat , sun) " +
		"VALUES" +
		"?,?,?,?,?,?,?,?"

	_, dbError := transaction.Exec(query, alarmID, repeatingIDs.Mon, repeatingIDs.Tue, repeatingIDs.Wed, repeatingIDs.Thu, repeatingIDs.Fri, repeatingIDs.Sat, repeatingIDs.Sun)
	return dbError
}

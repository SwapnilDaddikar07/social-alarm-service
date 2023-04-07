// Code generated by MockGen. DO NOT EDIT.
// Source: repository/alarm_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	constants "social-alarm-service/constants"
	db_model "social-alarm-service/db_model"
	transaction_manager "social-alarm-service/repository/transaction_manager"
	time "time"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockAlarmRepository is a mock of AlarmRepository interface.
type MockAlarmRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAlarmRepositoryMockRecorder
}

// MockAlarmRepositoryMockRecorder is the mock recorder for MockAlarmRepository.
type MockAlarmRepositoryMockRecorder struct {
	mock *MockAlarmRepository
}

// NewMockAlarmRepository creates a new mock instance.
func NewMockAlarmRepository(ctrl *gomock.Controller) *MockAlarmRepository {
	mock := &MockAlarmRepository{ctrl: ctrl}
	mock.recorder = &MockAlarmRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAlarmRepository) EXPECT() *MockAlarmRepositoryMockRecorder {
	return m.recorder
}

// CreateAlarmMetadata mocks base method.
func (m *MockAlarmRepository) CreateAlarmMetadata(ctx *gin.Context, transaction transaction_manager.Transaction, alarmId, userId string, alarmStartDateTime time.Time, alarmType constants.AlarmVisibility, description string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAlarmMetadata", ctx, transaction, alarmId, userId, alarmStartDateTime, alarmType, description)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAlarmMetadata indicates an expected call of CreateAlarmMetadata.
func (mr *MockAlarmRepositoryMockRecorder) CreateAlarmMetadata(ctx, transaction, alarmId, userId, alarmStartDateTime, alarmType, description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAlarmMetadata", reflect.TypeOf((*MockAlarmRepository)(nil).CreateAlarmMetadata), ctx, transaction, alarmId, userId, alarmStartDateTime, alarmType, description)
}

// GetAlarmMetadata mocks base method.
func (m *MockAlarmRepository) GetAlarmMetadata(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmMetadata", ctx, alarmId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmMetadata indicates an expected call of GetAlarmMetadata.
func (mr *MockAlarmRepositoryMockRecorder) GetAlarmMetadata(ctx, alarmId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmMetadata", reflect.TypeOf((*MockAlarmRepository)(nil).GetAlarmMetadata), ctx, alarmId)
}

// GetAllNonRepeatingAlarms mocks base method.
func (m *MockAlarmRepository) GetAllNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNonRepeatingAlarms", ctx, userId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNonRepeatingAlarms indicates an expected call of GetAllNonRepeatingAlarms.
func (mr *MockAlarmRepositoryMockRecorder) GetAllNonRepeatingAlarms(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNonRepeatingAlarms", reflect.TypeOf((*MockAlarmRepository)(nil).GetAllNonRepeatingAlarms), ctx, userId)
}

// GetAllRepeatingAlarms mocks base method.
func (m *MockAlarmRepository) GetAllRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRepeatingAlarms", ctx, userId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRepeatingAlarms indicates an expected call of GetAllRepeatingAlarms.
func (mr *MockAlarmRepositoryMockRecorder) GetAllRepeatingAlarms(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRepeatingAlarms", reflect.TypeOf((*MockAlarmRepository)(nil).GetAllRepeatingAlarms), ctx, userId)
}

// GetNonRepeatingAlarm mocks base method.
func (m *MockAlarmRepository) GetNonRepeatingAlarm(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNonRepeatingAlarm", ctx, alarmId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNonRepeatingAlarm indicates an expected call of GetNonRepeatingAlarm.
func (mr *MockAlarmRepositoryMockRecorder) GetNonRepeatingAlarm(ctx, alarmId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNonRepeatingAlarm", reflect.TypeOf((*MockAlarmRepository)(nil).GetNonRepeatingAlarm), ctx, alarmId)
}

// GetPublicNonExpiredAlarms mocks base method.
func (m *MockAlarmRepository) GetPublicNonExpiredAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, []db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicNonExpiredAlarms", ctx, userId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].([]db_model.Alarms)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPublicNonExpiredAlarms indicates an expected call of GetPublicNonExpiredAlarms.
func (mr *MockAlarmRepositoryMockRecorder) GetPublicNonExpiredAlarms(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicNonExpiredAlarms", reflect.TypeOf((*MockAlarmRepository)(nil).GetPublicNonExpiredAlarms), ctx, userId)
}

// GetPublicNonExpiredNonRepeatingAlarms mocks base method.
func (m *MockAlarmRepository) GetPublicNonExpiredNonRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicNonExpiredNonRepeatingAlarms", ctx, userId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicNonExpiredNonRepeatingAlarms indicates an expected call of GetPublicNonExpiredNonRepeatingAlarms.
func (mr *MockAlarmRepositoryMockRecorder) GetPublicNonExpiredNonRepeatingAlarms(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicNonExpiredNonRepeatingAlarms", reflect.TypeOf((*MockAlarmRepository)(nil).GetPublicNonExpiredNonRepeatingAlarms), ctx, userId)
}

// GetPublicNonExpiredRepeatingAlarms mocks base method.
func (m *MockAlarmRepository) GetPublicNonExpiredRepeatingAlarms(ctx *gin.Context, userId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicNonExpiredRepeatingAlarms", ctx, userId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicNonExpiredRepeatingAlarms indicates an expected call of GetPublicNonExpiredRepeatingAlarms.
func (mr *MockAlarmRepositoryMockRecorder) GetPublicNonExpiredRepeatingAlarms(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicNonExpiredRepeatingAlarms", reflect.TypeOf((*MockAlarmRepository)(nil).GetPublicNonExpiredRepeatingAlarms), ctx, userId)
}

// GetRepeatingAlarm mocks base method.
func (m *MockAlarmRepository) GetRepeatingAlarm(ctx *gin.Context, alarmId string) ([]db_model.Alarms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepeatingAlarm", ctx, alarmId)
	ret0, _ := ret[0].([]db_model.Alarms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepeatingAlarm indicates an expected call of GetRepeatingAlarm.
func (mr *MockAlarmRepositoryMockRecorder) GetRepeatingAlarm(ctx, alarmId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepeatingAlarm", reflect.TypeOf((*MockAlarmRepository)(nil).GetRepeatingAlarm), ctx, alarmId)
}

// InsertNonRepeatingDeviceAlarmID mocks base method.
func (m *MockAlarmRepository) InsertNonRepeatingDeviceAlarmID(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, deviceAlarmID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNonRepeatingDeviceAlarmID", ctx, transaction, alarmID, deviceAlarmID)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNonRepeatingDeviceAlarmID indicates an expected call of InsertNonRepeatingDeviceAlarmID.
func (mr *MockAlarmRepositoryMockRecorder) InsertNonRepeatingDeviceAlarmID(ctx, transaction, alarmID, deviceAlarmID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNonRepeatingDeviceAlarmID", reflect.TypeOf((*MockAlarmRepository)(nil).InsertNonRepeatingDeviceAlarmID), ctx, transaction, alarmID, deviceAlarmID)
}

// InsertRepeatingDeviceAlarmIDs mocks base method.
func (m *MockAlarmRepository) InsertRepeatingDeviceAlarmIDs(ctx *gin.Context, transaction transaction_manager.Transaction, alarmID string, repeatingIDs db_model.RepeatingAlarmIDs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertRepeatingDeviceAlarmIDs", ctx, transaction, alarmID, repeatingIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertRepeatingDeviceAlarmIDs indicates an expected call of InsertRepeatingDeviceAlarmIDs.
func (mr *MockAlarmRepositoryMockRecorder) InsertRepeatingDeviceAlarmIDs(ctx, transaction, alarmID, repeatingIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertRepeatingDeviceAlarmIDs", reflect.TypeOf((*MockAlarmRepository)(nil).InsertRepeatingDeviceAlarmIDs), ctx, transaction, alarmID, repeatingIDs)
}

// UpdateAlarmStatus mocks base method.
func (m *MockAlarmRepository) UpdateAlarmStatus(ctx *gin.Context, alarmId string, status constants.AlarmStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAlarmStatus", ctx, alarmId, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAlarmStatus indicates an expected call of UpdateAlarmStatus.
func (mr *MockAlarmRepositoryMockRecorder) UpdateAlarmStatus(ctx, alarmId, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAlarmStatus", reflect.TypeOf((*MockAlarmRepository)(nil).UpdateAlarmStatus), ctx, alarmId, status)
}
package request_model

import "social-alarm-service/db_model"

// CreateAlarmRequest TODO Check if we need to change int to int64. If notification package from flutter supports 64 , change BE type to int64.
type CreateAlarmRequest struct {
	AlarmStartDateTime        string                  `json:"alarm_start_date_time" binding:"required"`
	Private                   bool                    `json:"private"`
	Description               string                  `json:"description" binding:"max=100"`
	UserId                    string                  `json:"user_id" binding:"required"` // change this after moving ID to token.
	RepeatingDeviceAlarmIds   RepeatingDeviceAlarmIds `json:"repeating_device_alarm_ids"`
	NonRepeatingDeviceAlarmId *int                    `json:"non_repeating_device_alarm_id"`
}

type RepeatingDeviceAlarmIds struct {
	Mon *int `json:"mon"`
	Tue *int `json:"tue"`
	Wed *int `json:"wed"`
	Thu *int `json:"thu"`
	Fri *int `json:"fri"`
	Sat *int `json:"sat"`
	Sun *int `json:"sun"`
}

func (r RepeatingDeviceAlarmIds) ContainsAtleastOneRepeatingAlarm() bool {
	return r.Mon != nil ||
		r.Tue != nil ||
		r.Wed != nil ||
		r.Thu != nil ||
		r.Fri != nil ||
		r.Sat != nil ||
		r.Sun != nil
}

func (r RepeatingDeviceAlarmIds) MapToDBModel() db_model.RepeatingAlarmIDs {
	dbModel := db_model.RepeatingAlarmIDs{}
	if r.Mon != nil {
		dbModel.Mon = *r.Mon
	}
	if r.Tue != nil {
		dbModel.Tue = *r.Tue
	}
	if r.Wed != nil {
		dbModel.Wed = *r.Wed
	}
	if r.Thu != nil {
		dbModel.Thu = *r.Thu
	}
	if r.Fri != nil {
		dbModel.Fri = *r.Fri
	}
	if r.Sat != nil {
		dbModel.Sat = *r.Sat
	}
	if r.Sun != nil {
		dbModel.Sun = *r.Sun
	}
	return dbModel
}

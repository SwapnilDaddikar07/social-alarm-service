package constants

type AlarmVisibility string

const (
	AlarmPublicVisibility  AlarmVisibility = "PUBLIC"
	AlarmPrivateVisibility AlarmVisibility = "PRIVATE"
)

type AlarmStatus string

const (
	AlarmStatusON  AlarmStatus = "ON"
	AlarmStatusOFF AlarmStatus = "OFF"
)

func GetAlarmStatus(status string) AlarmStatus {
	if status == "ON" {
		return AlarmStatusON
	}
	return AlarmStatusOFF
}

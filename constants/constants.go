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

const (
	ContentTypeMP4   = "video/mp4"
	ContentTypeAudio = "audio/wave"
)

const MaxFileSizeInMB = 15

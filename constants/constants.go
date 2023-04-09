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

func (as AlarmStatus) String() string {
	return string(as)
}

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

const DateTimeLayout = "2006-01-02T15:04:05"

const MaxFileSizeInMB = 15

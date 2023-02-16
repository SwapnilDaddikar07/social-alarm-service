package response_model

import "social-alarm-service/db_model"

type MediaForAlarm struct {
	DisplayName string `json:"display_name"`
	MediaURL    string `json:"resource_url"`
}

func MapToMediaForAlarmResponseList(alarmMedia []db_model.GetMediaForAlarm) []MediaForAlarm {
	mediaForAlarm := make([]MediaForAlarm, 0)

	for _, entry := range alarmMedia {
		mediaForAlarm = append(mediaForAlarm, MediaForAlarm{
			DisplayName: entry.DisplayName,
			MediaURL:    entry.MediaURL,
		})
	}

	return mediaForAlarm
}

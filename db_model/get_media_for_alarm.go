package db_model

type GetMediaForAlarm struct {
	DisplayName string `db:"display_name"`
	MediaURL    string `db:"resource_url"`
}

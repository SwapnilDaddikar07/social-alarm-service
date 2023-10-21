package db_model

type User struct {
	UserId          string `db:"user_id"`
	PhoneNumber     string `db:"phone_number"`
	DisplayName     string `db:"display_name"`
	CurrentTimezone string `db:"current_timezone"`
}

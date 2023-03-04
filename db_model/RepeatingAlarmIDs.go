package db_model

type RepeatingAlarmIDs struct {
	Mon int `default:"-1"`
	Tue int `default:"-1"`
	Wed int `default:"-1"`
	Thu int `default:"-1"`
	Fri int `default:"-1"`
	Sat int `default:"-1"`
	Sun int `default:"-1"`
}

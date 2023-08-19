package utils

import (
	"github.com/nyaruka/phonenumbers"
	error2 "social-alarm-service/error"
)

func GetRegionCode(mobileNumber string) (string, *error2.ASError) {
	p, err := phonenumbers.Parse(mobileNumber, "")
	if err != nil {
		return "", error2.InternalServerError("invalid mobile number , cannot generate region code")
	}

	return phonenumbers.GetRegionCodeForCountryCode(int(p.GetCountryCode())), nil
}

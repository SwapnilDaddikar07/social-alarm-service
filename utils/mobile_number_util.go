package utils

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	error2 "social-alarm-service/error"
	"strconv"
)

func GetRegionCode(mobileNumber string) (string, *error2.ASError) {
	p, err := phonenumbers.Parse(mobileNumber, "")
	if err != nil {
		return "", error2.InternalServerError("invalid mobile number , cannot generate region code")
	}

	return phonenumbers.GetRegionCodeForCountryCode(int(p.GetCountryCode())), nil
}

func FormatNumbers(defaultRegionCode string, nonformattedMobileNumbers []string) []string {
	formattedMobileNumbers := make([]string, 0)

	for _, mobileNum := range nonformattedMobileNumbers {
		p, err := phonenumbers.Parse(mobileNum, defaultRegionCode)
		if err != nil {
			fmt.Println(fmt.Sprintf("could not parse the number %s ", mobileNum))
			continue
		}

		formattedNumber := fmt.Sprintf("%s%s%s", "+", fmt.Sprint(p.GetCountryCode()), strconv.FormatUint(p.GetNationalNumber(), 10))
		fmt.Println(fmt.Sprintf("unformatted number is %s and formatted number is %s", mobileNum, formattedNumber))

		formattedMobileNumbers = append(formattedMobileNumbers, formattedNumber)
	}

	return formattedMobileNumbers
}

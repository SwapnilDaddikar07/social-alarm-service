package error

import (
	"net/http"
)

type ASError struct {
	ErrorCode      string `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
	HttpStatusCode int    `json:"-"`
}

func InternalServerError(errorMessage string) *ASError {
	return &ASError{
		ErrorCode:      "ERR_INTERNAL_SERVER_ERROR",
		ErrorMessage:   errorMessage,
		HttpStatusCode: http.StatusInternalServerError,
	}
}

func BadRequestError(errorMessage string) *ASError {
	return &ASError{
		ErrorCode:      "ERR_BAD_REQUEST",
		ErrorMessage:   errorMessage,
		HttpStatusCode: http.StatusBadRequest,
	}
}

var AlarmIdMissing = &ASError{
	ErrorCode:      "ERR_ALARM_ID_MISSING",
	ErrorMessage:   "either repeating or non repeating alarm id should be provided",
	HttpStatusCode: http.StatusBadRequest,
}

var InvalidUserIdError = &ASError{
	ErrorCode:      "ERR_INVALID_USER_ID",
	ErrorMessage:   "user id does not exist",
	HttpStatusCode: http.StatusUnauthorized,
}

var DescriptionTooLongError = &ASError{
	ErrorCode:      "ERR_DESCRIPTION_TOO_LONG",
	ErrorMessage:   "description is too long.",
	HttpStatusCode: http.StatusBadRequest,
}

var InvalidAlarmTypeError = &ASError{
	ErrorCode:      "ERR_INVALID_ALARM_TYPE_ERROR",
	ErrorMessage:   "alarm cannot be repeating as well as non-repeating",
	HttpStatusCode: http.StatusBadRequest,
}

var InvalidAlarmDateTimeFormat = &ASError{
	ErrorCode:      "ERR_INVALID_ALARM_START_DATE_TIME_FORMAT",
	ErrorMessage:   "date format is incorrect",
	HttpStatusCode: http.StatusBadRequest,
}

var ContentTypeNotSupported = &ASError{
	ErrorCode:      "ERR_CONTENT_TYPE_NOT_SUPPORTED",
	ErrorMessage:   "content type is not supported",
	HttpStatusCode: http.StatusBadRequest,
}

var AlarmNotEligibleForMedia = &ASError{
	ErrorCode:      "ERR_ALARM_NOT_ELIGIBLE_FOR_MEDIA",
	ErrorMessage:   "alarm not eligible for media",
	HttpStatusCode: http.StatusUnauthorized,
}

var InvalidAlarmId = &ASError{
	ErrorCode:      "ERR_INVALID_ALARM_ID",
	ErrorMessage:   "alarm id is invalid",
	HttpStatusCode: http.StatusBadRequest,
}

var OperationNotAllowed = &ASError{
	ErrorCode:      "ERR_OPERATION_NOT_ALLOWED",
	ErrorMessage:   "you are not authorized to perform this operation",
	HttpStatusCode: http.StatusUnauthorized,
}

var NoPhoneNumbersInRequest = &ASError{
	ErrorCode:      "ERR_NO_PHONE_NUMBERS_IN_REQUEST",
	ErrorMessage:   "add contacts in your phone to load your dashboard",
	HttpStatusCode: http.StatusBadRequest,
}

var InvalidAlarmStatus = &ASError{
	ErrorCode:      "ERR_INVALID_ALARM_STATUS",
	ErrorMessage:   "status can only be on of off",
	HttpStatusCode: http.StatusBadRequest,
}

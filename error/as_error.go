package error

import "net/http"

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

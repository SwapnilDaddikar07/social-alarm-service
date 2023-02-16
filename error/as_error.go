package error

type ASError struct {
	ErrorCode      string `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
	HttpStatusCode int    `json:"-"`
}

func InternalServerError(errorMessage string) *ASError {
	return &ASError{
		ErrorCode:      "ERR_INTERNAL_SERVER_ERROR",
		ErrorMessage:   errorMessage,
		HttpStatusCode: 500,
	}
}

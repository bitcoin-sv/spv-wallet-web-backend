package api

// ErrorResponse is a struct that contains error message.
// This should be used for all error responses in http server.
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewResponseFromError creates ErrorResponse from error.
func NewErrorResponseFromError(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}

// NewResponseFromString creates ErrorResponse from string.
func NewErrorResponseFromString(err string) *ErrorResponse {
	return &ErrorResponse{
		Error: err,
	}
}

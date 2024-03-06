package errors

import "net/http"

type DomainError struct {
	Status  int    `json:"status"`  // http status mapping
	Code    string `json:"code"`    // domain error code
	Message string `json:"message"` // domain error message
}

func (err *DomainError) Error() string {
	return err.Message
}

var (
	DomainValidationError = &DomainError{
		Status:  http.StatusBadRequest,
		Code:    "2",
		Message: "Validation error",
	}
)

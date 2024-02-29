package errors

import "net/http"

type DomainError struct {
	Status  int    // http status mapping
	Code    string // domain error code
	Message string // domain error message
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

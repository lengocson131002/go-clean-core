package http

import (
	"errors"

	dErrors "github.com/lengocson131002/go-clean/pkg/errors"
)

type Result struct {
	Status  int         `json:"-"`       // Http status code
	Code    string      `json:"code"`    // Error code
	Message string      `json:"message"` // Message
	Details interface{} `json:"details"` //
}

type Response[T any] struct {
	Result Result `json:"result"` // Result
	Data   T      `json:"data"`   // Data
}

var (
	DefaultSuccessResponse = Response[interface{}]{
		Result: Result{
			Status:  200,
			Code:    "0",
			Message: "Success",
		},
		Data: nil,
	}

	DefaultFailureResponse = Response[interface{}]{
		Result: Result{
			Status:  500,
			Code:    "1",
			Message: "Internal Server Error",
		},
	}
)

func SuccessResponse[T any](data T) Response[T] {
	defRes := DefaultSuccessResponse
	return Response[T]{
		Result: defRes.Result,
		Data:   data,
	}
}

func FailureResponse(err error) Response[interface{}] {
	fRes := DefaultFailureResponse
	fRes.Result.Message = err.Error()

	var businessErr *dErrors.DomainError
	if errors.As(err, &businessErr) {
		fRes.Result = Result{
			Status:  businessErr.Status,
			Code:    businessErr.Code,
			Message: businessErr.Message,
		}
	}

	return fRes
}

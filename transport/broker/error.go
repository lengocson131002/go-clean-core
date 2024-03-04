package broker

import (
	"fmt"
	"time"
)

type EmptyRequestError struct{}

func (e EmptyRequestError) Error() string {
	return fmt.Sprintf("Empty broker request")
}

type InvalidDataFormatError struct{}

func (e InvalidDataFormatError) Error() string {
	return fmt.Sprintf("Invalid data format")
}

type RequestTimeoutResponse struct {
	Timeout time.Duration
}

func (e RequestTimeoutResponse) Error() string {
	return fmt.Sprintf("Request timeout exceeded. Timeout: %vs", e.Timeout.Seconds())
}

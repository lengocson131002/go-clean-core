package errors

import (
	"fmt"
)

var (
	InternalServerError = fmt.Errorf("Internal Server Error")
)

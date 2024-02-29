package broker

import "fmt"

type EmptyRequestError struct{}

func (e *EmptyRequestError) Error() string {
	return fmt.Sprintf("Empty broker request")
}

package es

import "fmt"

type ElasticSearchError struct {
	Status string `json:"status"`
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func (e ElasticSearchError) Error() string {
	return fmt.Sprintf("error from elastic search [%s] %s %s", e.Status, e.Type, e.Reason)
}

package healthchecks

import "sync"

type Integration struct {
	Name         string `json:"name"`
	Status       bool   `json:"status"`
	ResponseTime int64  `json:"response_time"`
	URL          string `json:"url,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ApplicationHealthDetailed struct {
	Name         string        `json:"name,omitempty"`
	Status       bool          `json:"status"`
	Version      string        `json:"version,omitempty"`
	Date         string        `json:"date"`
	Duration     int64         `json:"duration"`
	Integrations []Integration `json:"integration,omitempty"`
}

type HealthCheckHandler interface {
	Check(name string) Integration
}

type HealthChecker interface {
	AddLivenessCheck(name string, check HealthCheckHandler)
	AddReadinessCheck(name string, check HealthCheckHandler)
	LivenessCheck() ApplicationHealthDetailed
	RedinessCheck() ApplicationHealthDetailed
}

type HealthCheckerApplication struct {
	Name              string
	Version           string
	checksMutex       sync.RWMutex
	livenessCheckers  map[string]HealthCheckHandler
	readinessCheckers map[string]HealthCheckHandler
}

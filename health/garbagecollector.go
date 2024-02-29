package healthchecks

import (
	"fmt"
	"runtime"
	"time"
)

const (
	DEFAULT_GC_PAUSE_THRESHOLD = time.Duration(10) * time.Millisecond
)

type GarbageCollectionMaxChecker struct {
	threshold time.Duration
}

// GCMaxPauseCheck returns a Check that fails if any recent Go garbage
// collection pause exceeds the provided threshold.
func NewGarbageCollectionMaxChecker(threshold time.Duration) *GarbageCollectionMaxChecker {
	if threshold.Milliseconds() == 0 {
		threshold = DEFAULT_GC_PAUSE_THRESHOLD
	}

	return &GarbageCollectionMaxChecker{
		threshold: threshold,
	}
}

// Check implements HealthCheckHandler.
func (gc *GarbageCollectionMaxChecker) Check(name string) Integration {
	var (
		status       = true
		errorMessage = ""
		start        = time.Now()
	)

	thresholdNanos := uint64(gc.threshold.Nanoseconds())
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	for _, pause := range stats.PauseNs {
		if pause > thresholdNanos {
			status = false
			errorMessage = fmt.Sprintf("recent GC cycle took %s > %s", time.Duration(pause), gc.threshold)
			break
		}
	}

	return Integration{
		Status:       status,
		Name:         name,
		ResponseTime: time.Since(start).Milliseconds(),
		Error:        errorMessage,
	}

}

var _ HealthCheckHandler = (*GarbageCollectionMaxChecker)(nil)

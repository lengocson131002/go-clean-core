package healthchecks

import (
	"sync"
	"time"
)

func NewHealthChecker(name, version string) HealthChecker {
	app := &HealthCheckerApplication{
		livenessCheckers:  make(map[string]HealthCheckHandler),
		readinessCheckers: make(map[string]HealthCheckHandler),
		Name:              name,
		Version:           version,
	}

	return app
}

func (app *HealthCheckerApplication) runChecks(checks map[string]HealthCheckHandler) ApplicationHealthDetailed {
	var (
		start     = time.Now()
		wg        sync.WaitGroup
		checklist = make(chan Integration, len(checks))
		result    = ApplicationHealthDetailed{
			Name:         app.Name,
			Version:      app.Version,
			Status:       true,
			Date:         start.Format(time.RFC3339),
			Duration:     0,
			Integrations: []Integration{},
		}
	)

	wg.Add(len(checks))
	for name, handler := range checks {
		go func(name string, handler HealthCheckHandler) {
			checklist <- handler.Check(name)
			wg.Done()
		}(name, handler)
	}

	go func() {
		wg.Wait()
		close(checklist)
		result.Duration = time.Since(start).Milliseconds()
	}()

	for chk := range checklist {
		if !chk.Status {
			result.Status = false
		}
		result.Integrations = append(result.Integrations, chk)
	}

	return result
}

// AddLivenessCheck implements HealthChecker.
func (app *HealthCheckerApplication) AddLivenessCheck(name string, check HealthCheckHandler) {
	app.checksMutex.Lock()
	defer app.checksMutex.Unlock()
	app.livenessCheckers[name] = check
}

// AddReadinessCheck implements HealthChecker.
func (app *HealthCheckerApplication) AddReadinessCheck(name string, check HealthCheckHandler) {
	app.checksMutex.Lock()
	defer app.checksMutex.Unlock()
	app.readinessCheckers[name] = check
}

// LivenessCheck implements HealthChecker.
func (app *HealthCheckerApplication) LivenessCheck() ApplicationHealthDetailed {
	return app.runChecks(app.livenessCheckers)
}

// RedinessCheck implements HealthChecker.
func (app *HealthCheckerApplication) RedinessCheck() ApplicationHealthDetailed {
	return app.runChecks(app.readinessCheckers)
}

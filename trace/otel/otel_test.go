package otel

import (
	"context"
	"net/http"
	"testing"
)

func TestHttpTrace(t *testing.T) {
	var (
		ctx         = context.Background()
		serviceName = "test-service"
		method      = "GET"
		endpoint    = "http://google.com"
		exporter    = "localhost:4318"
	)

	client := http.Client{}

	tracer, err := NewOpenTelemetryTracer(ctx, serviceName, exporter)
	if err != nil {
		t.Fail()
	}

	ctx, finish := tracer.StartHttpClientTrace(
		ctx,
		"Ping to Google",
	)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)

	defer finish(ctx)

	client.Do(req)
	if err != nil {
		t.Fail()
	}

}

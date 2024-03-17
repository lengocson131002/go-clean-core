package otel

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/lengocson131002/go-clean-core/trace"
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

	tracer, err := NewOpenTelemetryTracer(
		ctx,
		trace.WithTraceServiceName(serviceName),
		trace.WithTraceExporterEndpoint(exporter))

	if err != nil {
		t.Fail()
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	ctx, finish := tracer.StartHttpClientTrace(
		ctx,
		"Ping to Google",
		trace.WithHttpRequest(req),
	)

	client.Do(req)
	if err != nil {
		t.Fail()
	}
	finish(ctx)

	time.Sleep(5 * time.Second)

}

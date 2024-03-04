package otel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lengocson131002/go-clean/pkg/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type openTelemetryTracer struct {
	serviceName      string
	exporterEndpoint string
}

func NewOpenTelemetryTracer(ctx context.Context, serviceName, exporterEndpoint string) (trace.Tracer, error) {
	// set global config trace
	if err := SetGlobalTracer(ctx, serviceName, exporterEndpoint); err != nil {
		return nil, err
	}

	return &openTelemetryTracer{
		serviceName:      serviceName,
		exporterEndpoint: exporterEndpoint,
	}, nil
}

// ExtractSpanInfo implements trace.Tracer.
func (*openTelemetryTracer) ExtractSpanInfo(ctx context.Context) *trace.SpanInfo {
	if span := oteltrace.SpanFromContext(ctx); span != nil {
		var traceID, spanID string

		if span.SpanContext().HasTraceID() {
			traceID = span.SpanContext().TraceID().String()
		}

		if span.SpanContext().HasSpanID() {
			spanID = span.SpanContext().SpanID().String()
		}

		return &trace.SpanInfo{
			TraceID: traceID,
			SpanID:  spanID,
		}
	}

	return nil
}

// HttpClientTrace implements trace.Tracer.
func (*openTelemetryTracer) StartHttpClientTrace(ctx context.Context, spanName string, opts ...trace.HttpClientTraceOption) (context.Context, trace.HttpClientTraceFinishFunc) {
	options := trace.HttpClientTraceOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	var (
		request          = options.Request
		endpoint, method string
	)

	if request != nil {
		endpoint = request.URL.String()
		method = request.Method

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))
	}

	tr := otel.Tracer(trace.HTTP_CLIENT)

	if spanName == "" {
		spanName = "http-client-call"
	}
	ctx, span := tr.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	span.SetAttributes(
		semconv.HTTPURLKey.String(endpoint),
		semconv.HTTPMethodKey.String(method),
	)

	return ctx, func(ctx context.Context, opts ...trace.HttpClientTraceFinishOption) {
		options := trace.HttpClientTraceFinishOptions{}
		for _, opt := range opts {
			opt(&options)
		}

		if span := oteltrace.SpanFromContext(ctx); span != nil {
			if res := options.Response; res != nil {
				span.SetAttributes(
					semconv.HTTPResponseContentLengthKey.Int64(res.ContentLength),
					semconv.HTTPStatusCodeKey.Int(res.StatusCode),
				)
			}

			span.End()
		}
	}
}

// StartDatabaseTrace implements trace.Tracer.
func (*openTelemetryTracer) StartDatabaseTrace(ctx context.Context, spanName string, opts ...trace.DatabaseTraceOption) (context.Context, trace.DatabaseTraceFinishFunc) {
	options := trace.DatabaseTraceOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	tr := otel.Tracer(trace.DATABASE)

	if spanName == "" {
		spanName = "database-operation"
	}
	ctx, span := tr.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))

	var (
		dbName      = options.DBName
		dbTableName = options.DBTableName
		sql         = options.DBSql
	)

	span.SetAttributes(
		attribute.String("db.name", dbName),
		attribute.String("db.table", dbTableName),
		attribute.String("db.sql", sql),
	)

	return ctx, func(ctx context.Context, opts ...trace.DatabaseTraceFinishOption) {
		options := trace.DatabaseTraceFinishOptions{}
		for _, opt := range opts {
			opt(&options)
		}

		if span := oteltrace.SpanFromContext(ctx); span != nil {
			span.End()
		}
	}
}

// StartGrpcClientTrace implements trace.Tracer.
func (*openTelemetryTracer) StartGrpcClientTrace(ctx context.Context, spanName string, opts ...trace.GrpcTraceOption) (context.Context, trace.GrpcTraceFinishFunc) {
	options := trace.GrpcTraceOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	tr := otel.Tracer(trace.GRPC_CLIENT)

	if spanName == "" {
		spanName = "grpc-client-call"
	}
	ctx, _ = tr.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))

	return ctx, func(ctx context.Context, opts ...trace.GrpcTraceFinishOption) {
		options := trace.GrpcTraceFinishOptions{}
		for _, opt := range opts {
			opt(&options)
		}
		if span := oteltrace.SpanFromContext(ctx); span != nil {
			if options.Error != nil {
				span.SetStatus(codes.Error, fmt.Sprintf("Error: %v", options.Error))
			}

			span.End()
		}
	}
}

// StartExternalTrace implements trace.Tracer.
func (*openTelemetryTracer) StartExternalTrace(ctx context.Context, spanName string, opts ...trace.ExternalTraceOption) (context.Context, trace.ExternalTraceFinishFunc) {
	options := trace.ExternalTraceOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	tr := otel.Tracer(trace.EXTERNAL)
	if spanName == "" {
		spanName = "external-client-call"
	}
	ctx, _ = tr.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))

	return ctx, func(ctx context.Context, opts ...trace.ExternalTraceFinishOption) {
		options := trace.ExternalTraceFinishOptions{}
		for _, opt := range opts {
			opt(&options)
		}
		if span := oteltrace.SpanFromContext(ctx); span != nil {
			if options.Response != nil {
				resJson, _ := json.Marshal(options.Response)
				span.SetAttributes(attribute.String("external.response", string(resJson)))
			}

			span.End()
		}
	}
}

// StartInternalTrace implements trace.Tracer.
func (*openTelemetryTracer) StartInternalTrace(ctx context.Context, spanName string, opts ...trace.InternalTraceOption) (context.Context, trace.InternalTraceFinishFunc) {
	options := trace.InternalTraceOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	tr := otel.Tracer(trace.INTERNAL)

	if spanName == "" {
		spanName = "internal-operation"
	}
	ctx, span := tr.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	if req := options.Request; req != nil {
		reqJson, _ := json.Marshal(req)
		span.SetAttributes(attribute.String("internal.request", string(reqJson)))
	}

	return ctx, func(ctx context.Context, opts ...trace.InternalTraceFinishOption) {
		options := trace.InternalTraceFinishOptions{}
		for _, opt := range opts {
			opt(&options)
		}

		if span := oteltrace.SpanFromContext(ctx); span != nil {
			if res := options.Response; res != nil {
				resJson, _ := json.Marshal(res)
				span.SetAttributes(attribute.String("internal.response", string(resJson)))
			}
			span.End()
		}
	}
}

// Setup global tracing configurations
func SetGlobalTracer(ctx context.Context, serviceName, exporterEndpoint string) error {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(exporterEndpoint))

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

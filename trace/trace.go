package trace

import "context"

type SpanFinishFunc func(context.Context, interface{})

type SpanInfo struct {
	TraceID string
	SpanID  string
}

type Tracer interface {
	// Extract span info from the context
	ExtractSpanInfo(context.Context) *SpanInfo

	// Used for tracing HTTP Client Call
	StartHttpClientTrace(ctx context.Context, spanName string, opts ...HttpClientTraceOption) (context.Context, HttpClientTraceFinishFunc)

	// Used for tracing GRPC Client Call
	StartGrpcClientTrace(ctx context.Context, spanName string, opts ...GrpcTraceOption) (context.Context, GrpcTraceFinishFunc)

	// Used for tracing DATABASE call
	StartDatabaseTrace(ctx context.Context, spanName string, opts ...DatabaseTraceOption) (context.Context, DatabaseTraceFinishFunc)

	// Used for tracing other external services (queue, cache,...)
	StartExternalTrace(ctx context.Context, spanName string, opts ...ExternalTraceOption) (context.Context, ExternalTraceFinishFunc)

	// Used for tracing interal functions
	StartInternalTrace(ctx context.Context, spanName string, opts ...InternalTraceOption) (context.Context, InternalTraceFinishFunc)
}

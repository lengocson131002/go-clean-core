package trace

import "context"

type SpanFinishFunc func(context.Context, interface{})

type SpanInfo struct {
	TraceID        string `json:"traceID"`
	SpanID         string `json:"spanID"`
	ClientIP       string `json:"clientIP"`
	HttpMethod     string `json:"httpMethod"`
	ServiceDomain  string `json:"serviceDomain"`
	OperatorName   string `json:"operatorName"`
	StepName       string `json:"stepName"`
	UserAgent      string `json:"userAgent"`
	User           string `json:"user"`
	ProcessTime    string `json:"processTime"`
	RemoteHost     string `json:"remoteHost"`
	XForwardedFor  string `json:"xForwardedFor"`
	ContentLength  string `json:"contentLength"`
	StatusResponse string `json:"statusResponse"`
}

/*
*
-	Log format: %d{yyyy-MM-dd} %d{HH:mm:ss.SSS} %level [%thread] %logger{50} [%X{X-B3-TraceId},%X{X-B3-SpanId}] [%X{systemTraceId}] [%X{clientIP}] [%X{httpMethod}] [%X{serviceDomain}] [%X{operatorName}] [%X{stepName}] [%X{req.userAgent}] [%X{user}] [%X{processTime}] [%X{req.remoteHost}] [%X{req.xForwardedFor}] [%X{contentLength}] [%X{statusResponse}]
*/
type Tracer interface {
	// Extract span info from the context
	ExtractSpanInfo(context.Context) SpanInfo

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

package trace

type TraceOptions struct {
	ServiceName      string
	ExporterEndpoint string
}

type TraceOption func(*TraceOptions)

func WithTraceServiceName(serviceName string) TraceOption {
	return func(options *TraceOptions) {
		options.ServiceName = serviceName
	}
}

func WithTraceExporterEndpoint(ep string) TraceOption {
	return func(options *TraceOptions) {
		options.ExporterEndpoint = ep
	}
}

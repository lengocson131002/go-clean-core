package trace

import "context"

// External trace options
type ExternalTraceOption func(*ExternalTraceOptions)

type ExternalTraceOptions struct {
	DestServiceName string
	Endpoint        string
	Request         interface{}
}

func WithExternalServiceName(serviceName string) ExternalTraceOption {
	return func(opts *ExternalTraceOptions) {
		opts.DestServiceName = serviceName
	}
}

func WithExternalEndpoint(endpoint string) ExternalTraceOption {
	return func(opts *ExternalTraceOptions) {
		opts.Endpoint = endpoint
	}
}

type ExternalTraceFinishOption func(*ExternalTraceFinishOptions)

type ExternalTraceFinishOptions struct {
	Response interface{}
}

func WithExternalResponse(res interface{}) ExternalTraceFinishOption {
	return func(opts *ExternalTraceFinishOptions) {
		opts.Response = res
	}
}

type ExternalTraceFinishFunc func(context.Context, ...ExternalTraceFinishOption)

package trace

import "context"

// GRPC trace options
type GrpcTraceOption func(*GrpcTraceOptions)

type GrpcTraceOptions struct{}

type GrpcTraceFinishOption func(*GrpcTraceFinishOptions)

type GrpcTraceFinishOptions struct {
	Error error
}

func WithGrpcError(err error) GrpcTraceFinishOption {
	return func(opts *GrpcTraceFinishOptions) {
		opts.Error = err
	}
}

type GrpcTraceFinishFunc func(context.Context, ...GrpcTraceFinishOption)

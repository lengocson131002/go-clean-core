package trace

import "context"

// Internal trace options
type InternalTraceOption func(*InternalTraceOptions)

type InternalTraceOptions struct {
	Request interface{}
}

func WithInternalRequest(request interface{}) InternalTraceOption {
	return func(options *InternalTraceOptions) {
		options.Request = request
	}
}

type InternalTraceFinishOption func(*InternalTraceFinishOptions)

type InternalTraceFinishOptions struct {
	Response interface{}
	Error    error
}

func WithInternalResponse(response interface{}) InternalTraceFinishOption {
	return func(options *InternalTraceFinishOptions) {
		options.Response = response
	}
}

func WithErrorResponse(err error) InternalTraceFinishOption {
	return func(options *InternalTraceFinishOptions) {
		options.Error = err
	}
}

type InternalTraceFinishFunc func(context.Context, ...InternalTraceFinishOption)

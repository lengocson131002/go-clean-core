package trace

import (
	"context"
	"net/http"
)

// HTTP trace options
type HttpClientTraceOption func(*HttpClientTraceOptions)

type HttpClientTraceOptions struct {
	Request *http.Request
}

func WithHttpRequest(request *http.Request) HttpClientTraceOption {
	return func(options *HttpClientTraceOptions) {
		options.Request = request
	}
}

type HttpClientTraceFinishOption func(*HttpClientTraceFinishOptions)

type HttpClientTraceFinishOptions struct {
	Response *http.Response
}

func WithHttpResponse(response *http.Response) HttpClientTraceFinishOption {
	return func(options *HttpClientTraceFinishOptions) {
		options.Response = response
	}
}

type HttpClientTraceFinishFunc func(context.Context, ...HttpClientTraceFinishOption)

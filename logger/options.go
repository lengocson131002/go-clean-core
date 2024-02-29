package logger

import (
	"context"
	"io"

	"github.com/lengocson131002/go-clean-core/trace"
)

type Option func(*Options)

type Options struct {
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Out io.Writer
	// fields to always be logged
	Fields map[string]interface{}
	// Caller skip frame count for file:line info
	CallerSkipCount int
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// Alternative options
	Context context.Context
	//tracer for extract Trace ID, Span ID
	Tracer trace.Tracer
}

// WithFields set default fields for the logger.
func WithFields(fields map[string]interface{}) Option {
	return func(args *Options) {
		args.Fields = fields
	}
}

// WithLevel set default level for the logger.
func WithLevel(level Level) Option {
	return func(args *Options) {
		args.Level = level
	}
}

// WithOutput set default output writer for the logger.
func WithOutput(out io.Writer) Option {
	return func(args *Options) {
		args.Out = out
	}
}

// WithCallerSkipCount set frame count to skip.
func WithCallerSkipCount(c int) Option {
	return func(args *Options) {
		args.CallerSkipCount = c
	}
}

func WithTracer(t trace.Tracer) Option {
	return func(args *Options) {
		args.Tracer = t
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

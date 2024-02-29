package logger

import "context"

type Logger interface {
	// Init initializes options
	Init(options ...Option) error
	// The Logger options
	Options() Options
	// Fields set fields to always be logged
	Fields(fields map[string]interface{}) Logger

	// Log writes a log entry
	Log(ctx context.Context, level Level, args ...interface{})
	// Logf writes a formatted log entry
	Logf(ctx context.Context, level Level, format string, args ...interface{})
	// Log Trace
	Trace(ctx context.Context, args ...interface{})
	// Logf Trace
	Tracef(ctx context.Context, format string, args ...interface{})
	// Log Debug
	Debug(ctx context.Context, args ...interface{})
	// Logf Debug
	Debugf(ctx context.Context, format string, args ...interface{})
	// Log Info
	Info(ctx context.Context, args ...interface{})
	// Logf Info
	Infof(ctx context.Context, format string, args ...interface{})
	// Log Warn
	Warn(ctx context.Context, args ...interface{})
	// Logf Warn
	Warnf(ctx context.Context, format string, args ...interface{})
	// Log Error
	Error(ctx context.Context, args ...interface{})
	// Logf Error
	Errorf(ctx context.Context, format string, args ...interface{})
	// Log Fatal
	Fatal(ctx context.Context, args ...interface{})
	// Logf Fatal
	Fatalf(ctx context.Context, format string, args ...interface{})

	// String returns the name of logger
	String() string
}

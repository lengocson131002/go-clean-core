package logrus

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/sirupsen/logrus"
)

type entryLogger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry

	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type logrusLogger struct {
	Logger entryLogger
	opts   LogrusOptions
}

func (l *logrusLogger) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if formatter, ok := l.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		l.opts.Formatter = formatter
	}

	if hs, ok := l.opts.Context.Value(hooksKey{}).(logrus.LevelHooks); ok {
		l.opts.Hooks = hs
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}

	switch ll := l.opts.Context.Value(logrusLoggerKey{}).(type) {
	case *logrus.Logger:
		// overwrite default options
		l.opts.Level = logrusToLoggerLevel(ll.GetLevel())
		l.opts.Out = ll.Out
		l.opts.Formatter = ll.Formatter
		l.opts.Hooks = ll.Hooks
		l.opts.ReportCaller = ll.ReportCaller
		l.opts.ExitFunc = ll.ExitFunc
		l.Logger = ll
	case *logrus.Entry:
		// overwrite default options
		el := ll.Logger
		l.opts.Level = logrusToLoggerLevel(el.GetLevel())
		l.opts.Out = el.Out
		l.opts.Formatter = el.Formatter
		l.opts.Hooks = el.Hooks
		l.opts.ReportCaller = el.ReportCaller
		l.opts.ExitFunc = el.ExitFunc
		l.Logger = ll
	case nil:
		log := logrus.New() // defaults
		log.SetLevel(loggerToLogrusLevel(l.opts.Level))
		log.SetOutput(l.opts.Out)
		log.SetFormatter(l.opts.Formatter)
		log.ReplaceHooks(l.opts.Hooks)
		log.SetReportCaller(l.opts.ReportCaller)
		log.ExitFunc = l.opts.ExitFunc
		l.Logger = log
	default:
		return fmt.Errorf("invalid logrus type: %T", ll)
	}

	return nil
}

func (l *logrusLogger) String() string {
	return "logrus"
}

func (l *logrusLogger) Fields(fields map[string]interface{}) logger.Logger {
	return &logrusLogger{l.Logger.WithFields(fields), l.opts}
}

func (l *logrusLogger) Log(ctx context.Context, level logger.Level, args ...interface{}) {
	var entry = l.getLogEntry(ctx)
	entry.Log(loggerToLogrusLevel(level), args...)

}

func (l *logrusLogger) Logf(ctx context.Context, level logger.Level, format string, args ...interface{}) {
	var entry = l.getLogEntry(ctx)
	entry.Logf(loggerToLogrusLevel(level), format, args...)
}

func (l *logrusLogger) getLogEntry(ctx context.Context) entryLogger {
	var entry = l.Logger
	entry = entry.WithContext(ctx)
	return entry
}

func (l *logrusLogger) Options() logger.Options {
	return l.opts.Options
}

// Debug implements logger.Logger.
func (l *logrusLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.DebugLevel, args...)
}

// Debugf implements logger.Logger.
func (l *logrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.DebugLevel, format, args...)
}

// Error implements logger.Logger.
func (l *logrusLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.ErrorLevel, args...)
}

// Errorf implements logger.Logger.
func (l *logrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.ErrorLevel, format, args...)
}

// Fatal implements logger.Logger.
func (l *logrusLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.FatalLevel, args...)
}

// Fatalf implements logger.Logger.
func (l *logrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.FatalLevel, format, args...)
}

// Info implements logger.Logger.
func (l *logrusLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.InfoLevel, args...)
}

// Infof implements logger.Logger.
func (l *logrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.InfoLevel, format, args...)
}

// Trace implements logger.Logger.
func (l *logrusLogger) Trace(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.TraceLevel, args...)
}

// Tracef implements logger.Logger.
func (l *logrusLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.TraceLevel, format, args...)
}

// Warn implements logger.Logger.
func (l *logrusLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.WarnLevel, args...)
}

// Warnf implements logger.Logger.
func (l *logrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.WarnLevel, format, args...)
}

// New builds a new logger based on options.
func NewLogrusLogger(opts ...logger.Option) logger.Logger {
	// Default options
	var l logrusLogger

	loggerOpts := logger.Options{
		MaskSensitiveData: true,
		Level:             logger.InfoLevel,
		Fields:            make(map[string]interface{}),
		Out:               os.Stderr,
		Context:           context.Background(),
		CallerSkipCount:   7,
	}
	l.opts = LogrusOptions{
		Options: loggerOpts,
		Formatter: &LoggingFormatter{
			logger: &l,
		},
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}
	_ = l.Init(opts...)
	return &l
}

func loggerToLogrusLevel(level logger.Level) logrus.Level {
	switch level {
	case logger.TraceLevel:
		return logrus.TraceLevel
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	case logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func logrusToLoggerLevel(level logrus.Level) logger.Level {
	switch level {
	case logrus.TraceLevel:
		return logger.TraceLevel
	case logrus.DebugLevel:
		return logger.DebugLevel
	case logrus.InfoLevel:
		return logger.InfoLevel
	case logrus.WarnLevel:
		return logger.WarnLevel
	case logrus.ErrorLevel:
		return logger.ErrorLevel
	case logrus.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}

type LoggingFormatter struct {
	logger logger.Logger
}

func (l LoggingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var (
		opts = l.logger.Options()
	)
	// Get the file and line number where the log was called
	_, filename, line, _ := runtime.Caller(opts.CallerSkipCount)

	// Get the script name from the full file path
	scriptName := filepath.Base(filename)

	// Common format
	message := fmt.Sprintf("%s %s %s [%s:%d] %s\n",
		entry.Time.Format("2006-01-02"),       // Date
		entry.Time.Format("15:04:05.123"),     // Time
		strings.ToUpper(entry.Level.String()), // Log level
		scriptName,                            //Script name
		line,                                  // Line number
		entry.Message,                         // Log message
	)

	// Log format: %d{yyyy-MM-dd} %d{HH:mm:ss.SSS} %level [%thread] %logger{50} [%X{X-B3-TraceId},%X{X-B3-SpanId}] [%X{systemTraceId}] [%X{clientIP}] [%X{httpMethod}] [%X{serviceDomain}] [%X{operatorName}] [%X{stepName}] [%X{req.userAgent}] [%X{user}] [%X{processTime}] [%X{req.remoteHost}] [%X{req.xForwardedFor}] [%X{contentLength}] [%X{statusResponse}]
	if traceEnabled := opts.Tracer != nil; traceEnabled {
		spanInfo := opts.Tracer.ExtractSpanInfo(entry.Context)
		message = fmt.Sprintf("%s %s %s [%s] [%s] [%s:%d] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] %s\n",
			entry.Time.Format("2006-01-02"),
			entry.Time.Format("15:04:05.123"),
			strings.ToUpper(entry.Level.String()),
			spanInfo.TraceID,
			spanInfo.SpanID,
			scriptName,
			line,
			spanInfo.ClientIP,
			spanInfo.HttpMethod,
			spanInfo.ServiceDomain,
			spanInfo.OperatorName,
			spanInfo.StepName,
			spanInfo.UserAgent,
			spanInfo.User,
			spanInfo.ProcessTime,
			spanInfo.RemoteHost,
			spanInfo.XForwardedFor,
			spanInfo.ContentLength,
			spanInfo.StatusResponse,
			entry.Message,
		)
	}

	if opts.MaskSensitiveData {
		message = logger.MaskSensitiveData(message)
	}

	return []byte(message), nil
}

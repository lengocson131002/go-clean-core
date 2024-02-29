package trace

import "context"

// Database trace options
type DatabaseTraceOption func(*DatabaseTraceOptions)

type DatabaseTraceOptions struct {
	DBName      string
	DBTableName string
	DBSql       string
}

func WithDBName(dbName string) DatabaseTraceOption {
	return func(options *DatabaseTraceOptions) {
		options.DBName = dbName
	}
}

func WithDBTableName(tableName string) DatabaseTraceOption {
	return func(options *DatabaseTraceOptions) {
		options.DBTableName = tableName
	}
}

func WithDBSql(sql string) DatabaseTraceOption {
	return func(options *DatabaseTraceOptions) {
		options.DBSql = sql
	}
}

type DatabaseTraceFinishOption func(*DatabaseTraceFinishOptions)

type DatabaseTraceFinishOptions struct {
	Error error
}

type DatabaseTraceFinishFunc func(context.Context, ...DatabaseTraceFinishOption)

func WithDBError(err error) DatabaseTraceFinishOption {
	return func(options *DatabaseTraceFinishOptions) {
		options.Error = err
	}
}

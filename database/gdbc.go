package database

import (
	"context"
	"database/sql"
)

// SqlGdbc (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type SqlGdbc interface {
	Transactor
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// Used this in repositories
type Gdbc struct {
	sDB SqlGdbc
}

// Exec implements SqlGdbc.
func (g *Gdbc) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return g.getConnection(ctx).Exec(query, args...)
}

// Get implements SqlGdbc.
func (g *Gdbc) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return g.getConnection(ctx).Get(dest, query, args...)
}

// Prepare implements SqlGdbc.
func (g *Gdbc) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return g.getConnection(ctx).Prepare(query)
}

// Query implements SqlGdbc.
func (g *Gdbc) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return g.getConnection(ctx).Query(query, args...)
}

// QueryRow implements SqlGdbc.
func (g *Gdbc) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return g.getConnection(ctx).QueryRow(query, args...)
}

// Select implements SqlGdbc.
func (g *Gdbc) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return g.getConnection(ctx).Select(dest, query, args...)
}

func (g *Gdbc) getConnection(ctx context.Context) SqlGdbc {
	s := ExtractTx(ctx)
	if s != nil {
		return s
	}
	return g.sDB
}

func (g *Gdbc) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return g.sDB.WithinTransaction(ctx, txFunc)
}

func (g *Gdbc) WithinTransactionOptions(ctx context.Context, txFunc func(ctx context.Context) error, txOptions *sql.TxOptions) error {
	return g.sDB.WithinTransactionOptions(ctx, txFunc, txOptions)
}

package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// SqlxDBx is the sqlx.DB based implementation of GDBC
type SqlxDBTx struct {
	DB *sqlx.DB
}

// SqlxConnx is the sqlx.Tx based implementation of GDBC
type SqlxConnTx struct {
	DB *sqlx.Tx
}

func NewSqlxDBGdbc(db *sqlx.DB) *SqlxDBTx {
	return &SqlxDBTx{db}
}

func NewSqlxConnGdbc(db *sqlx.Tx) *SqlxConnTx {
	return &SqlxConnTx{db}
}

// Get implements database.SqlGdbc.
func (s *SqlxDBTx) Get(dest interface{}, query string, args ...interface{}) error {
	return s.DB.Get(dest, query, args...)
}

// Select implements database.SqlGdbc.
func (s *SqlxDBTx) Select(dest interface{}, query string, args ...interface{}) error {
	return s.DB.Select(dest, query, args)
}

// Exec implements SqlGdbc.
func (s *SqlxDBTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.Exec(query, args...)
}

// Prepare implements SqlGdbc.
func (s *SqlxDBTx) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.Prepare(query)
}

// Query implements SqlGdbc.
func (s *SqlxDBTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.Query(query, args)
}

// QueryRow implements SqlGdbc.
func (s *SqlxDBTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.DB.QueryRow(query, args...)
}

// Exec implements SqlGdbc.
func (s *SqlxConnTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.Exec(query, args...)
}

// Prepare implements SqlGdbc.
func (s *SqlxConnTx) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.Prepare(query)
}

// Query implements SqlGdbc.
func (s *SqlxConnTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.Query(query, args)
}

// QueryRow implements SqlGdbc.
func (s *SqlxConnTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.DB.QueryRow(query, args...)
}

// Get implements database.SqlGdbc.
func (s *SqlxConnTx) Get(dest interface{}, query string, args ...interface{}) error {
	return s.DB.Get(dest, query, args...)
}

// Select implements database.SqlGdbc.
func (s *SqlxConnTx) Select(dest interface{}, query string, args ...interface{}) error {
	return s.DB.Select(dest, query, args)
}

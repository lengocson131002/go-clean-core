package database

import (
	"github.com/jmoiron/sqlx"
)

type SqlxDatabaseConnector struct {
}

func NewSqlxDatabaseConnector() DatabaseConnector {
	return &SqlxDatabaseConnector{}
}

func (c *SqlxDatabaseConnector) Connect(driverName string, dsn string, poolOptions *PoolOptions) (*Gdbc, error) {
	db := sqlx.MustConnect(driverName, dsn)

	err := db.Ping()
	if err != nil {
		if db != nil {
			err = db.Close()
		}
		return nil, err
	}

	if poolOptions != nil {
		db.SetMaxIdleConns(poolOptions.MaxIdleCount)
		db.SetMaxOpenConns(poolOptions.MaxOpen)
		db.SetConnMaxLifetime(poolOptions.MaxLifetime)
		db.SetConnMaxIdleTime(poolOptions.MaxIdleTime)
	}

	return &Gdbc{
		&SqlxDBTx{db},
	}, err
}

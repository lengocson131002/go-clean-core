package sqlx

import (
	"github.com/jmoiron/sqlx"
	"github.com/lengocson131002/go-clean-core/database"
)

type SqlxDatabase struct {
}

func NewSqlxDatabaseConnector() database.Database {
	return &SqlxDatabase{}
}

func (c *SqlxDatabase) Connect(driverName string, dsn string, poolOptions *database.PoolOptions) (*database.Gdbc, error) {
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

	return &database.Gdbc{
		Executor: &SqlxDBTx{
			DB: db,
		},
	}, err
}

package database

type DatabaseConnector interface {
	// Connect to database using:
	// 1. drivername
	// 2. dsn(connection string)
	// 3. poolOptions: nil if no need to configure pool
	Connect(drivername string, dsn string, poolOptions *PoolOptions) (*Gdbc, error)
}

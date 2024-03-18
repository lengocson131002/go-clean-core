#### Usage

```go

type YugabyteConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	SslMode               string
	IdleConnection        int
	MaxConnection         int
	MaxLifeTimeConnection int //seconds
	MaxIdleTimeConnection int // seconds
}

func GetYugabyteConfig(cfg config.Configure) *YugabyteConfig {
	username := cfg.GetString("DB_YUGABYTE_USER")
	password := cfg.GetString("DB_YUGABYTE_PASSWORD")
	host := cfg.GetString("DB_YUGABYTE_HOST")
	port := cfg.GetInt("DB_YUGABYTE_PORT")
	sslmode := "disable"
	database := cfg.GetString("DB_YUGABYTE_DBNAME")
	idleConnection := cfg.GetInt("DB_YUGABYTE_POOL_IDLE_CONNECTION")
	maxConnection := cfg.GetInt("DB_YUGABYTE_MAX_POOL_SIZE")
	maxLifeTimeConnection := cfg.GetInt("DB_YUGABYTE_MAX_LIFE_TIME")
	maxLifeIdleConnection := cfg.GetInt("DB_YUGABYTE_IDLE_TIMEOUT")

	return &YugabyteConfig{
		Username:              username,
		Password:              password,
		Host:                  host,
		Port:                  port,
		Database:              database,
		SslMode:               sslmode,
		IdleConnection:        idleConnection,
		MaxConnection:         maxConnection,
		MaxLifeTimeConnection: maxLifeTimeConnection,
		MaxIdleTimeConnection: maxLifeIdleConnection,
	}
}

func GetDatabaseConnector() database.Database {
	return sqlx.NewSqlxDatabaseConnector()
}

func GetMasterDataDatabase(y *YugabyteConfig, conn database.Database) *data.MasterDataDatabase {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s", y.Username, y.Password, y.Host, y.Port, y.Database, fmt.Sprintf("sslmode=%s", y.SslMode))

	db, err := conn.Connect("postgres", dsn, &database.PoolOptions{
		MaxIdleCount: y.IdleConnection,
		MaxOpen:      y.MaxConnection,
		MaxLifetime:  time.Duration(y.MaxLifeTimeConnection) * time.Second,
		MaxIdleTime:  time.Duration(y.MaxIdleTimeConnection) * time.Second,
	})

	if err != nil {
		panic(err)
	}

	return &data.MasterDataDatabase{
		DB: db,
	}
}

```

``` go
package data

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean/pkg/database"
	"github.com/lengocson131002/go-clean/pkg/trace"
)

// Interface for metarepository
type MasterDataRepository interface {
	GetTemplateRequest(ctx context.Context, templateName string) (string, error)
}

type templateEntity struct {
	tName     string `db:"template_name"`
	tRequest  string `db:"template_request"`
	tResponse string `db:"template_response"`
}

type MasterDataDatabase struct {
	DB *database.Gdbc
}

type masterDataRepository struct {
	db     *database.Gdbc
	tracer trace.Tracer
}

func NewMasterDataRepository(db *MasterDataDatabase, tracer trace.Tracer) MasterDataRepository {
	return &masterDataRepository{
		db:     db.DB,
		tracer: tracer,
	}
}

func (repo *masterDataRepository) GetTemplateRequest(ctx context.Context, templateName string) (string, error) {
	sql := "SELECT TEMPLATE_NAME, TEMPLATE_REQUEST, TEMPLATE_RESPONSE FROM GW_XSLTEMPLATES WHERE TEMPLATE_NAME = $1"

	ctx, finish := repo.tracer.StartDatabaseTrace(
		ctx,
		"get template request from master database",
		trace.WithDBTableName("GW_XSLTEMPLATES"),
		trace.WithDBSql(sql),
	)

	defer finish(ctx)

	row := repo.db.QueryRow(ctx, sql, templateName)
	if row == nil {
		return "", fmt.Errorf("Template not found")
	}

	if row.Err() != nil {
		return "", fmt.Errorf("failed to get template: %v", row.Err())
	}

	tempEntity := new(templateEntity)

	row.Scan(&tempEntity.tName, &tempEntity.tRequest, &tempEntity.tResponse)

	return tempEntity.tRequest, nil
}

```


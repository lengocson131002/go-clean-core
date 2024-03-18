#### Usage

``` go

func GetConfigure() config.Configure {
	var file viper.ConfigFile = ".env"
	cfg, err := viper.NewViperConfig(&file)
	if err != nil {
		panic(err)
	}

	return cfg
}

```

```go
cfg, err := GetConfigure()
if err != nil {
    panic(err)
}

username := cfg.GetString("DB_YUGABYTE_USER")
password := cfg.GetString("DB_YUGABYTE_PASSWORD")
host := cfg.GetString("DB_YUGABYTE_HOST")
port := cfg.GetInt("DB_YUGABYTE_PORT")

```
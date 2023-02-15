package sql

type sqlConfig struct {
	Dialect            string
	Host               string
	Port               string
	Username           string
	Password           string
	DbName             string
	ConnectionString   string
	MaxIdleConnections *int
	MaxOpenConnections *int
	MaxLifetimeSeconds *int
}

type config struct {
	Debug bool
	Sql   sqlConfig
}

package sql

type config struct {
	Debug              bool
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
	MaxIdleTimeSeconds *int
}

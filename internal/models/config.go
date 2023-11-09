package models

var DefaultConfig = Config{}

type Config struct {
	Postgres  PostgresConfig
	Webserver WebserverConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type WebserverConfig struct {
	Addr       string
	PublicAddr string
	TLS        TLSConfig
}

type TLSConfig struct {
	Enabled bool
	Cert    string
	Key     string
}
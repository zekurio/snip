package models

import "github.com/zekurio/snip/pkg/randutils"

var DefaultConfig = Config{
	Webserver: WebserverConfig{
		Addr:       ":80",
		PublicAddr: "http://localhost:80",
		AccessToken: AccessToken{
			Secret:          randutils.ForceRandBase64Str(64),
			LifetimeSeconds: 10 * 60,
		},
		TLS: TLSConfig{
			Enabled: false,
			Cert:    "",
			Key:     "",
		},
	},
}

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
	Addr        string
	PublicAddr  string
	TLS         TLSConfig
	AccessToken AccessToken
}

type TLSConfig struct {
	Enabled bool
	Cert    string
	Key     string
}

type AccessToken struct {
	Secret          string
	LifetimeSeconds int
}

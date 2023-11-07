package webserver

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

package config

var DefaultConfig = Config{
	Addr: "localhost:8080",
	Redirects: []Redirects{
		{
			ID:   "test",
			Link: "https://google.com",
		},
	},
}

type Config struct {
	Addr      string      `json:"addr"`
	Redirects []Redirects `json:"redirects"`
}

type Redirects struct {
	ID   string `json:"id"`   // ID of te redirect, basically the part after the domain
	Link string `json:"link"` // The link to redirect to
}

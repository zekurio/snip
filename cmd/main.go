package main

import (
	"flag"
	"net/http"

	"github.com/zekurio/redirector/pkg/config"
)

var (
	flagConfigPath = flag.String("c", "config.toml", "Path to config file")
)

func main() {
	// TODO load config with Parse from pkg/config/parser.go
	cfg, err := config.Parse(*flagConfigPath, "REDIRECTOR_", config.DefaultConfig)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		for _, redirect := range cfg.Redirects {
			if path == redirect.ID {
				http.Redirect(w, r, redirect.Link, http.StatusFound)
				return
			}
		}
		http.NotFound(w, r)
	})

	http.ListenAndServe(cfg.Addr, nil)
}

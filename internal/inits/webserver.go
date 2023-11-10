package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/webserver"
	"github.com/zekurio/snip/internal/util/static"
)

func InitWebserver(ctn di.Container) (ws *webserver.WebServer) {
	cfg := ctn.Get(static.DiConfig).(models.Config)

	ws, err := webserver.New(ctn)
	logrus.Info("Initlializing webserver")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create webserver")
		return
	}

	go func() {
		if err = ws.ListenAndServeBlocking(); err != nil {
			logrus.WithError(err).Fatal("Failed to start webserver")
		}
	}()
	logrus.Info("Webserver started",
		"bindAddr", cfg.Webserver.Addr, "publicAddr", cfg.Webserver.PublicAddr)

	return
}

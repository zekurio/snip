package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/services/webserver"
)

func InitWebserver(ctn di.Container) (ws *webserver.WebServer) {
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

	return
}

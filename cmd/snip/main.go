package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/redirector/internal/services/config"
	"github.com/zekurio/redirector/internal/services/util/static"
)

var (
	flagConfigPath = flag.String("c", "config.toml", "Path to config file")
)

func main() {
	// Parse command line flags
	flag.Parse()

	diBuilder, err := di.NewBuilder()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create DI builder")
	}

	// Config dependency
	diBuilder.Add(di.Def{
		Name: static.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Parse(*flagConfigPath, "SNIP_", config.DefaultConfig)
		},
	})

	// Webserver dependency
	diBuilder.Add(di.Def{
		Name: static.DiWebserver,
		Build: func(ctn di.Container) (interface{}, error) {
			return nil, nil
		},
	})

	// Build dependency injection container
	ctn := diBuilder.Build()

	// Tear down dependency instances
	defer func(ctn di.Container) {
		err := ctn.DeleteWithSubContainers()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to tear down dependency instances")
		}
	}(ctn)

	// Block main go routine until one of the following
	// specified exit sys calls occure.
	logrus.Info("Started event loop. Stop with CTRL-C...")

	logrus.Info("Initialization finished")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

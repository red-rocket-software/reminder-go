package main

import (
	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/server"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()

	app := server.New(logger)

	logger.Debugf("Starting server on port %s", cfg.HTTP.Port)

	if err := app.Run(cfg.HTTP.IP, cfg.HTTP.Port); err != nil {
		logger.Fatalf("%s", err.Error())
	}
}

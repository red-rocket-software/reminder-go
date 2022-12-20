package main

import (
	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/pkg/utils"
	"github.com/red-rocket-software/reminder-go/server"
	log "github.com/sirupsen/logrus"
)

func init() {
	utils.ConfigureLogger()
}

func main() {
	cfg := config.GetConfig()
	log.Info(cfg.HTTP.IP)

	app := server.New()

	log.Info("Starting server on port 8080")

	if err := app.Run("localhost", "8080"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}

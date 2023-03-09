package main

import (
	"context"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/server"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
)

//	@title			Reminder App API
//	@version		1.0
//	@description	API Server for Reminder Application

//	@host		localhost:8000
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth
func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Getting new db client...")
	postgresClient, err := postgresql.NewClient(ctx, 5, *cfg)
	if err != nil {
		logger.Fatalf("Error create new db client:%v\n", err)
	}
	defer postgresClient.Close()

	todoStorage := storage.NewStorageTodo(postgresClient, &logger)

	app := server.New(ctx, logger, todoStorage, *cfg)
	logger.Debugf("Starting server on port %s", cfg.HTTP.Port)

	if err := app.Run(cfg); err != nil {
		logger.Fatalf("%s", err.Error())
	}
}

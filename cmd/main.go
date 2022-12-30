package main

import (
	"context"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/storage/todo"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
	"github.com/red-rocket-software/reminder-go/server"
)

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

	todo.NewStorageTodo(postgresClient, &logger)

	app := server.New(logger)
	logger.Debugf("Starting server on port %s", cfg.HTTP.Port)

	if err := app.Run(cfg.HTTP.IP, cfg.HTTP.Port); err != nil {
		logger.Fatalf("%s", err.Error())
	}
}

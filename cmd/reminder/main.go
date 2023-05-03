package main

import (
	"context"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/reminder/server"
	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
	"google.golang.org/api/option"
)

//	@title			Reminder App API
//	@version		1.0
//	@description	API Server for Reminder Application

// @host		localhost:8000
// @BasePath	/
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

	// creating firebase client
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	fireClient, err := firestore.NewClient(ctx, opt)
	if err != nil {
		logger.Errorf("Failed to Auth a Firestore Client: %v", err)
		return
	}

	app := server.New(ctx, logger, todoStorage, fireClient, *cfg)
	logger.Debugf("Starting reminder server on port %s", cfg.HTTP.Port)

	if err := app.Run(cfg); err != nil {
		logger.Fatalf("%s", err.Error())
	}

}

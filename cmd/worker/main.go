package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
	"github.com/red-rocket-software/reminder-go/worker"
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
	defer postgresClient.Close()

	todoStorage := storage.NewStorageTodo(postgresClient, &logger)

	worker := worker.NewWorker(ctx, todoStorage, *cfg)

	//run worker in scheduler
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	stop := make(chan bool)

	ticker := time.NewTicker(time.Second * 5) // worker runs every 5 second

	go func() {
		defer func() { stop <- true }()
		for {
			select {
			case <-ticker.C:
				err = worker.Process()
				if err != nil {
					logger.Errorf("error to process worker: %v", err)
					return
				}
			case <-stop:
				logger.Info("closing goroutine")
				return
			}
		}

	}()
	<-c
	defer ticker.Stop()

	stop <- true

	<-stop
	logger.Info("Stop application")
}

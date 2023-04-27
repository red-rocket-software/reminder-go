package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/red-rocket-software/reminder-go/config"
	todoStorage "github.com/red-rocket-software/reminder-go/internal/reminder/storage"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
	"github.com/red-rocket-software/reminder-go/worker"
	"google.golang.org/api/option"
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

	// creating firebase client
	logger.Info("Getting new firebase client...")
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	fireClient, err := firestore.NewClient(ctx, opt)
	if err != nil {
		logger.Errorf("Failed to Auth a Firestore Client: %v", err)
		return
	}

	remindStorage := todoStorage.NewStorageTodo(postgresClient, &logger)

	newWorker := worker.NewWorker(ctx, remindStorage, fireClient, *cfg)

	//run worker in scheduler
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	stop := make(chan bool)

	ticker := time.NewTicker(time.Second * 10) // worker runs every 10 second

	go func() {
		defer func() { stop <- true }()
		for {
			select {
			case <-ticker.C:
				logger.Info("processing worker notification...")
				err = newWorker.ProcessSendNotification()
				if err != nil {
					logger.Errorf("error to process worker send notification: %v", err)
					return
				}
				logger.Info("processing worker deadline notification...")
				err = newWorker.ProcessSendDeadlineNotification()
				if err != nil {
					logger.Errorf("error to process worker send deadline notification: %v", err)
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

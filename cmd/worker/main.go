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
	"github.com/red-rocket-software/reminder-go/workers/notifier"
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

	newWorker := notifier.NewWorker(ctx, remindStorage, fireClient, *cfg)

	//run workers in scheduler
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	stop := make(chan error)

	ticker := time.NewTicker(time.Second * 10) // workers runs every 10 second

	go func() {
		for {
			select {
			case <-ticker.C:
				err = newWorker.ProcessSendNotification()
				if err != nil {
					logger.Errorf("error to process workers send notification: %v", err)
					stop <- err
				}
				err = newWorker.ProcessSendDeadlineNotification()
				if err != nil {
					logger.Errorf("error to process workers send deadline notification: %v", err)
					stop <- err
				}
			case <-stop:
				logger.Info("closing goroutine")
				return
			}
		}

	}()
	<-c
	defer ticker.Stop()

	<-stop
	logger.Info("Stop application")
}

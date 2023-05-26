package server

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/config"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

type Server struct {
	S              *http.Server
	Router         *mux.Router
	Logger         logging.Logger
	TodoStorage    model.TodoRepository
	ConfigsStorage model.ConfigRepository
	FireClient     firestore.Client
	ctx            context.Context
	config         config.Config
}

// New returns new Server.
func New(ctx context.Context, logger logging.Logger, todoStorage model.TodoRepository, configsStorage model.ConfigRepository, fireClient firestore.Client, cfg config.Config) *Server {
	server := &Server{
		ctx:            ctx,
		Logger:         logger,
		TodoStorage:    todoStorage,
		ConfigsStorage: configsStorage,
		FireClient:     fireClient,
		config:         cfg,
	}
	return server
}

// Run start server on IP address an PORT passed in parameters
func (server *Server) Run(cfg *config.Config) error {
	server.S = &http.Server{
		Addr:           ":" + cfg.HTTP.Port,
		Handler:        server.ConfigureReminderRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.S.ListenAndServe(); err != nil {
			server.Logger.Fatalf("Failed to listen and : %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	<-ctx.Done()

	server.Logger.Info("Shutting down")
	os.Exit(0)

	return server.S.Shutdown(ctx)
}

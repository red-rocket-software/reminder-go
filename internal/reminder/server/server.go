package server

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

type Server struct {
	S           *http.Server
	Router      *mux.Router
	Logger      logging.Logger
	TodoStorage storage.ReminderRepo
	FireClient  *auth.Client
	ctx         context.Context
	config      config.Config
}

// New returns new Server.
func New(ctx context.Context, logger logging.Logger, storage storage.ReminderRepo, fireClient *auth.Client, cfg config.Config) *Server {

	server := &Server{
		ctx:         ctx,
		Logger:      logger,
		TodoStorage: storage,
		FireClient:  fireClient,
		config:      cfg,
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

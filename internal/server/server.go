package server

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/pkg/logging"

	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	S           *http.Server
	Router      *mux.Router
	Logger      logging.Logger
	TodoStorage storage.ReminderRepo
	ctx         context.Context
}

// func New returns new Server. You should pass logger as a parameter
func New(ctx context.Context, logger logging.Logger, storage storage.ReminderRepo) *Server {

	server := &Server{Logger: logger, TodoStorage: storage, ctx: ctx}
	return server
}

// func Run start server on IP address an PORT passed in parameters
func (server *Server) Run(cfg *config.Config) error {

	server.S = &http.Server{

		Addr:           cfg.HTTP.IP + ":" + cfg.HTTP.Port,
		Handler:        server.ConfigureRouter(),
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	<-ctx.Done()

	server.Logger.Info("Shutting down")
	os.Exit(0)

	return server.S.Shutdown(ctx)
}

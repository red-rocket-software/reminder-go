package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/router"
	"github.com/red-rocket-software/reminder-go/pkg/logging"

	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	S      *http.Server
	Router *mux.Router
	Logger logging.Logger
}

// func New returns new Server. You should pass DB as a parameter
func New(logger logging.Logger) *Server {
	return &Server{Logger: logger}
}

// func Run start server on IP address an PORT passed in parameters
func (server *Server) Run(ip string, port string) error {

	server.S = &http.Server{

		Addr:           ip + ":" + port,
		Handler:        router.ConfigureRouter(server.Logger),
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

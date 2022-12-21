package server

import (
	"context"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	s      *http.Server
	logger logging.Logger
}

// func New returns new Server. You should pass DB as a parameter
func New(logger logging.Logger) *Server {
	return &Server{logger: logger}
}

// func Run start server on IP address an PORT passed in parameters
func (server *Server) Run(ip string, port string) error {

	server.s = &http.Server{

		Addr: ip + ":" + port,
		// add router here!
		// Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.s.ListenAndServe(); err != nil {
			server.logger.Fatalf("Failed to listen and : %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	<-ctx.Done()

	server.logger.Info("Shutting down")
	os.Exit(0)

	return server.s.Shutdown(ctx)
}

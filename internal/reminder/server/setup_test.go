package server

import (
	"context"
	"os"
	"testing"

	"github.com/red-rocket-software/reminder-go/config"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

func newTestServer(todoStorage model.TodoRepository, configsStorage model.ConfigRepository) *Server {
	logger := logging.GetLogger()
	cfg := config.Config{}

	server := New(context.Background(), logger, todoStorage, configsStorage, cfg)

	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

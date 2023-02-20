package server

import (
	"context"
	"os"
	"testing"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

func newTestServer(store storage.ReminderRepo) *Server {
	logger := logging.GetLogger()
	//cfg := config.GetConfig()
	cfg := config.Config{}

	server := New(context.Background(), logger, store, cfg)

	return server
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

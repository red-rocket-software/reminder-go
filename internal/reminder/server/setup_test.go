package server

import (
	"context"
	"os"
	"testing"

	"github.com/red-rocket-software/reminder-go/config"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"google.golang.org/api/option"
)

func newTestServer(todoStorage model.TodoRepository, configsStorage model.ConfigRepository) *Server {
	logger := logging.GetLogger()
	cfg := config.Config{}

	// creating firebase client
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	fireClient, _ := firestore.NewClient(context.Background(), opt)

	server := New(context.Background(), logger, todoStorage, configsStorage, fireClient, cfg)

	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

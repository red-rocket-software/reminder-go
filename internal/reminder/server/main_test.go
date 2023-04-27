package server

import (
	"context"
	"os"
	"testing"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"google.golang.org/api/option"
)

func newTestServer(store storage.ReminderRepo) *Server {
	logger := logging.GetLogger()
	//cfg := config.GetConfig()
	cfg := config.Config{}

	// creating firebase client
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	fireClient, _ := firestore.NewClient(context.Background(), opt)

	server := New(context.Background(), logger, store, fireClient, cfg)

	return server
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

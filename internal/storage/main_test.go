package storage

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

var testStorage *StorageTodo

func TestMain(m *testing.M) {
	//cfg := config.GetConfig()
	logger := logging.GetLogger()

	pClient, err := pgxpool.New(context.Background(), "postgres://root:secret@localhost:5432/reminder?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db...")
	}

	testStorage = NewStorageTodo(pClient, &logger)

	os.Exit(m.Run())
}

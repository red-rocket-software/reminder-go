package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

var (
	ErrNoNewMigrations = errors.New("no change")
)

var testStorage *TodoStorage

func TestMain(m *testing.M) {
	//cfg := config.GetConfig()
	logger := logging.GetLogger()

	pClient, err := pgxpool.New(context.Background(), "postgres://root:secret@localhost:5432/test_reminder?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db...")
	}

	testStorage = NewStorageTodo(pClient, &logger)

	os.Exit(m.Run())
}

func DropEverythingInDatabase() error {
	m, err := migrate.New("file://db/migration", "postgres://root:secret@localhost:5432/test_reminder?sslmode=disable")
	if err != nil {
		return err
	}

	if err := m.Drop(); err != nil {
		return errors.WithStack(err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil || dbErr != nil {
		return errors.Errorf("srcErr: %v and dbErr: %v", srcErr, dbErr)
	}

	fmt.Println("db drop successfully")

	return nil
}

func RunUpMigrations() error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../db/migrations")
	migrationDir := filepath.Join("file://" + basePath)
	m, err := migrate.New(migrationDir, "postgres://root:secret@localhost:5432/test_reminder?sslmode=disable")
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, ErrNoNewMigrations) {
			return errors.WithStack(err)
		}
	}
	m.Close()
	return nil
}

package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

var testTodoStorage model.TodoRepository
var testConfigStorage model.ConfigRepository
var pClient *pgxpool.Pool

func TestMain(m *testing.M) {
	logger := logging.GetLogger()
	var err error

	pClient, err = pgxpool.New(context.Background(), "postgres://root:secret@localhost:5432/test_coa?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db...")
	}

	testTodoStorage = NewStorageTodo(pClient, &logger)
	testConfigStorage = NewConfigsStorage(pClient, &logger)

	os.Exit(m.Run())
}

// SeedTodos seed todos for tests with notification
func SeedTodosForDeadline() ([]model.Todo, error) {
	date := time.Date(2023, time.April, 1, 1, 0, 0, 0, time.UTC)
	now := time.Now().Truncate(1 * time.Millisecond).UTC()
	finishedDate := time.Date(2023, time.April, 1, 2, 0, 0, 0, time.UTC)
	dateNotifyPeriod, _ := time.Parse(time.RFC3339, time.Now().Truncate(time.Minute).Format(time.RFC3339))

	userID, err := SeedUserConfig()
	if err != nil {
		fmt.Println(err)
		return []model.Todo{}, err
	}

	b := true

	todos := []model.Todo{
		{
			Description:    "tes1",
			Title:          "tes1",
			UserID:         userID,
			CreatedAt:      now,
			DeadlineAt:     date,
			Completed:      false,
			DeadlineNotify: &b,
			NotifyPeriod:   []time.Time{dateNotifyPeriod},
		},
		{
			Description: "tes2",
			Title:       "tes2",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
			Completed:   true,
			FinishedAt:  &finishedDate,
		},
		{
			Description: "tes3",
			Title:       "tes3",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
		},
		{
			Description: "tes4",
			Title:       "tes4",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
		},
		{
			Description: "tes5",
			Title:       "tes5",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
		},
	}

	for i := range todos {
		const sql = `INSERT INTO reminder.todo ("Description", "Title", "User", "CreatedAt", "DeadlineAt", "FinishedAt", "Completed", "DeadlineNotify", "NotifyPeriod") 
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning "ID"`

		row := pClient.QueryRow(context.Background(), sql, todos[i].Description, todos[i].Title, todos[i].UserID, todos[i].CreatedAt, todos[i].DeadlineAt, todos[i].FinishedAt, todos[i].Completed, todos[i].DeadlineNotify, todos[i].NotifyPeriod)

		err := row.Scan(&todos[i].ID)
		if err != nil {
			return []model.Todo{}, fmt.Errorf("Error create remind: %v", err)
		}

	}

	return todos, nil
}

// SeedTodos seed todos for tests
func SeedTodos() ([]model.Todo, error) {
	date := time.Date(2023, time.April, 1, 1, 0, 0, 0, time.UTC)
	now := time.Now().Truncate(1 * time.Millisecond).UTC()
	finishedDate := time.Date(2023, time.April, 1, 2, 0, 0, 0, time.UTC)

	userID, err := SeedUserConfig()
	if err != nil {
		fmt.Println(err)
		return []model.Todo{}, err
	}

	todos := []model.Todo{
		{
			Description: "tes1",
			Title:       "tes1",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
			Completed:   false,
		},
		{
			Description: "tes2",
			Title:       "tes2",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
			Completed:   true,
			FinishedAt:  &finishedDate,
		},
		{
			Description: "tes3",
			Title:       "tes3",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
		},
		{
			Description: "tes4",
			Title:       "tes4",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
		},
		{
			Description: "tes5",
			Title:       "tes5",
			UserID:      userID,
			CreatedAt:   now,
			DeadlineAt:  date,
		},
	}

	for i := range todos {
		const sql = `INSERT INTO reminder.todo ("Description", "Title", "User", "CreatedAt", "DeadlineAt", "FinishedAt", "Completed") 
				 VALUES ($1, $2, $3, $4, $5, $6, $7) returning "ID"`

		row := pClient.QueryRow(context.Background(), sql, todos[i].Description, todos[i].Title, todos[i].UserID, todos[i].CreatedAt, todos[i].DeadlineAt, todos[i].FinishedAt, todos[i].Completed)

		err := row.Scan(&todos[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error create remind: %v", err)
		}

	}

	return todos, nil
}

// Truncate removes all seed data from the test database.
func Truncate() error {
	stmt := "TRUNCATE TABLE reminder.todo, reminder.users_configs;"

	if _, err := pClient.Exec(context.Background(), stmt); err != nil {
		return fmt.Errorf("truncate test database tables %v", err)
	}

	return nil
}

// SeedUserConfig seed todos for tests
func SeedUserConfig() (string, error) {
	userConfig := model.UserConfigs{
		ID:           "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
		Notification: true,
		Period:       2,
		CreatedAt:    time.Now(),
	}

	var ID string

	const sql = `INSERT INTO reminder.users_configs ("ID", "Notification", "Period", "CreatedAt") 
				 VALUES ($1, $2, $3, $4) returning "ID"`

	row := pClient.QueryRow(context.Background(), sql, userConfig.ID, userConfig.Notification, userConfig.Period, userConfig.CreatedAt)

	err := row.Scan(&ID)
	if err != nil {
		return "", fmt.Errorf("error create userConfig: %v", err)
	}

	return ID, nil
}

package todo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

// StorageTodo handles database communication with PostgreSQL.
type StorageTodo struct {
	// Postgres database.PGX
	Postgres *pgxpool.Pool
	// Logrus logger
	logger *logging.Logger
}

func NewStorageTodo(postgres *pgxpool.Pool, logger *logging.Logger) *StorageTodo {
	return &StorageTodo{Postgres: postgres, logger: logger}
}

func (s *StorageTodo) GetAllReminds(ctx context.Context) ([]model.Todo, error) {
	return nil, nil
	//TODO implement me
}

// CreateRemind  create new remind in PostgreSQL
func (s *StorageTodo) CreateRemind(ctx context.Context, todo model.Todo) error {
	const sql = `INSERT INTO todo ("Id", "Description", "CreatedAt", "DeadlineAt", "FinishedAt") 
				 VALUES ($1, $2, $3, $4)`

	_, err := s.Postgres.Exec(ctx, sql, todo.ID, todo.Description, todo.CreatedAt, todo.DeadlineAt, todo.FinishedAt)
	if err != nil {
		s.logger.Errorf("Error create remind: %v", err)
		return err
	}
	return nil
}
func (s *StorageTodo) UpdateRemind(ctx context.Context, id string) (model.Todo, error) {
	return model.Todo{}, nil
	//TODO implement me
}
func (s *StorageTodo) DeleteRemind(ctx context.Context, id string) error {
	return nil
	//TODO implement me
}
func (s *StorageTodo) GetRemindByID(ctx context.Context, id string) (model.Todo, error) {
	return model.Todo{}, nil
	//TODO implement me
}
func (s *StorageTodo) GetComplitedReminds(ctx context.Context) ([]model.Todo, error) {
	return nil, nil
	//TODO implement me
}
func (s *StorageTodo) GetNewReminds(ctx context.Context) ([]model.Todo, error) {
	return nil, nil
	//TODO implement me
}

package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

// NewStorageTodo  return new SorageTodo with Postgres pool and logger
func NewStorageTodo(postgres *pgxpool.Pool, logger *logging.Logger) *StorageTodo {
	return &StorageTodo{Postgres: postgres, logger: logger}
}

func (s *StorageTodo) GetAllReminds(ctx context.Context) ([]model.Todo, error) {
	return nil, nil
	//TODO implement me
}

// CreateRemind  store new remind entity to DB PostgreSQL
func (s *StorageTodo) CreateRemind(ctx context.Context, todo model.Todo) error {
	var id int
	const sql = `INSERT INTO todo ("Description", "CreatedAt", "DeadlineAt") 
				 VALUES ($1, $2, $3) returning "Id"`
	row := s.Postgres.QueryRow(ctx, sql, todo.Description, todo.CreatedAt, todo.DeadlineAt)
	err := row.Scan(&id)
	if err != nil {
		s.logger.Errorf("Error create remind: %v", err)
		return err
	}
	return nil
}

// UpdateRemind update remind, can change Description, Completed and FihishedAt if Completed = true
func (s *StorageTodo) UpdateRemind(ctx context.Context, id int, input model.TodoUpdate) error {
	const sql = `UPDATE todo SET "Description" = $1, "FinishedAt" = $2, "Completed" = $3 WHERE "Id" = $4`

	ct, err := s.Postgres.Exec(ctx, sql, input.Description, input.FinishedAt, input.Completed, id)
	if err != nil {
		s.logger.Printf("unable to update remind %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("remind not found")
	}

	return nil
}
func (s *StorageTodo) DeleteRemind(ctx context.Context, id string) error {
	return nil
	//TODO implement me
}

// GetRemindByID takes out one remind from PostgreSQL by id
func (s *StorageTodo) GetRemindByID(ctx context.Context, id int) (model.Todo, error) {
	var todo model.Todo

	const sql = `SELECT "Id", "Description", "CreatedAt", "DeadlineAt", "Completed" FROM todo
    WHERE "Id" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, id)

	err := row.Scan(
		&todo.ID,
		&todo.Description,
		&todo.CreatedAt,
		&todo.DeadlineAt,
		&todo.Completed,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Todo{}, nil
	}
	if err != nil {
		s.logger.Printf("cannot get product from database: %v\n", err)
		return model.Todo{}, errors.New("cannot get product from database")
	}

	return todo, nil
}

func (s *StorageTodo) GetComplitedReminds(ctx context.Context) ([]model.Todo, error) {
	return nil, nil
	//TODO implement me
}
func (s *StorageTodo) GetNewReminds(ctx context.Context) ([]model.Todo, error) {
	return nil, nil
	//TODO implement me
}

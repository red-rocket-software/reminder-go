package storage

import (
	"context"
	"errors"
	"fmt"

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

type FetchParam struct {
	Limit    int
	CursorID int
}

// GetAllReminds return all todos in DB PostgreSQL
func (s *StorageTodo) GetAllReminds(ctx context.Context, fetchParams FetchParams) ([]model.Todo, int, error) {
	var reminds []model.Todo

	const sql = `SELECT "Id", "Description", "CreatedAt", "DeadlineAt", "FinishedAt", "Completed" FROM todo WHERE Id > $1  ORDER BY "CreatedAt" DESC LIMIT $2`

	rows, err := s.Postgres.Query(ctx, sql, fetchParams.Cursor, fetchParams.Limit)

	if err != nil {
		s.logger.Errorf("error get all reminds from db: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var remind model.Todo

		if err := rows.Scan(
			&remind.ID,
			&remind.Description,
			&remind.CreatedAt,
			&remind.DeadlineAt,
			&remind.FinishedAt,
			&remind.Completed,
		); err != nil {
			s.logger.Errorf("remind doesnt exist: %v", err)
			return nil, 0, err
		}
		reminds = append(reminds, remind)
	}

	var nextCursor int

	if len(reminds) > 0 {
		nextCursor = reminds[len(reminds)-1].ID
	}

	return reminds, nextCursor, nil
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

// DeleteRemind deletes remind from DB
func (s *StorageTodo) DeleteRemind(ctx context.Context, id int) error {
	const sql = `DELETE FROM todo WHERE id = $1`
	res, err := s.Postgres.Exec(ctx, sql, id)

	if err != nil {
		s.logger.Errorf("error don't found remind: %v", err)
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Errorf("Error delete remind: %v", err)
		return ErrDeleteFailed
	}

	return nil
}

// GetRemindByID takes out one remind from PostgreSQL by id
func (s *StorageTodo) GetRemindByID(ctx context.Context, id int) (model.Todo, error) {
	var todo model.Todo

	const sql = `SELECT "Id", "Description", "CreatedAt", "DeadlineAt", "Completed", "FinishedAt" FROM todo
    WHERE "Id" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, id)

	err := row.Scan(
		&todo.ID,
		&todo.Description,
		&todo.CreatedAt,
		&todo.DeadlineAt,
		&todo.Completed,
		&todo.FinishedAt,
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

// GetNewReminds get all no completed reminds from DB with pagination.
func (s *StorageTodo) GetNewReminds(ctx context.Context, params FetchParam) ([]model.Todo, int, error) {
	sql := `SELECT * FROM todo WHERE "Completed" = false`

	//if passed cursorID we add condition to query
	if params.CursorID > 0 {
		sql += fmt.Sprintf(` AND "Id" < %d`, params.CursorID)
	}

	//always add sort and LIMIT to query
	sql += fmt.Sprintf(` ORDER BY "CreatedAt" DESC LIMIT %d`, params.Limit)

	rows, err := s.Postgres.Query(ctx, sql)
	if err != nil {
		s.logger.Errorf("error to select completed reminds: %v", err)
		return nil, 0, err
	}

	var reminds []model.Todo

	for rows.Next() {
		var remind model.Todo

		if err := rows.Scan(&remind.ID,
			&remind.Description,
			&remind.CreatedAt,
			&remind.DeadlineAt,
			&remind.FinishedAt,
			&remind.Completed,
		); err != nil {
			s.logger.Error("remind doesn't exist %v", err)
			return nil, 0, err
		}

		reminds = append(reminds, remind)
	}

	var nextCursor int
	if len(reminds) > 0 {
		nextCursor = reminds[len(reminds)-1].ID
	}

	return reminds, nextCursor, nil
}

package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

var _ model.TodoRepository = (*TodoStorage)(nil)

// TodoStorage handles database communication with PostgreSQL.
type TodoStorage struct {
	// Postgres database.PGX
	Postgres *pgxpool.Pool
	// Logrus logger
	logger *logging.Logger
}

// NewStorageTodo  return new SorageTodo with Postgres pool and logger
func NewStorageTodo(postgres *pgxpool.Pool, logger *logging.Logger) model.TodoRepository {
	return &TodoStorage{Postgres: postgres, logger: logger}
}

// GetReminds return all todos in DB PostgreSQL
func (s *TodoStorage) GetReminds(ctx context.Context, params model.FetchParams, userID string) ([]model.Todo, int, int, error) {
	var reminds []model.Todo
	var sql string

	switch params.FilterByQuery {
	case "current":
		sql = fmt.Sprintf(`SELECT *, (
SELECT COUNT(*) FROM reminder.todo WHERE "User" = '%s' AND "Completed" = false) as total_count
FROM reminder.todo WHERE "User" = '%s' AND "Completed" = false`, userID, userID)
	case "completed":
		sql = fmt.Sprintf(`SELECT *, (
SELECT COUNT(*) FROM reminder.todo WHERE "User" = '%s' AND "Completed" = true) as total_count
FROM reminder.todo WHERE "User" = '%s' AND "Completed" = true`, userID, userID)
	case "all":
		sql = fmt.Sprintf(`SELECT *, (
SELECT COUNT(*) FROM reminder.todo WHERE "User" = '%s') as total_count
FROM reminder.todo WHERE "User" = '%s'`, userID, userID)
	default:
		return nil, 0, 0, errors.New("wrong filterParams value")
	}

	if params.Cursor > 0 {
		switch params.FilterByDate {
		case "DESC":
			sql += fmt.Sprintf(` AND "%s" < (SELECT "%s" FROM reminder.todo WHERE "ID" = %d)`, params.FilterByDate, params.FilterByDate, params.Cursor)
		case "ASC":
			sql += fmt.Sprintf(` AND "%s" > (SELECT "%s" FROM reminder.todo WHERE "ID" = %d)`, params.FilterByDate, params.FilterByDate, params.Cursor)
		}
	}

	if params.StartRange != "" {
		sql += fmt.Sprintf(` AND "FinishedAt" BETWEEN '%s' AND '%s'`, params.StartRange, params.EndRange)
	}

	sql += fmt.Sprintf(` ORDER BY "%s" %s LIMIT %d`, params.FilterByDate, params.FilterBySort, params.Limit)

	rows, err := s.Postgres.Query(ctx, sql)

	if err != nil {
		s.logger.Errorf("error get all reminds from db: %v", err)
		return []model.Todo{}, 0, 0, err
	}
	defer rows.Close()

	var totalCount int

	for rows.Next() {
		var remind model.Todo

		if err := rows.Scan(
			&remind.ID,
			&remind.UserID,
			&remind.Title,
			&remind.Description,
			&remind.CreatedAt,
			&remind.DeadlineAt,
			&remind.FinishedAt,
			&remind.Completed,
			&remind.Notificated,
			&remind.DeadlineNotify,
			&remind.NotifyPeriod,
			&totalCount,
		); err != nil {
			s.logger.Errorf("remind doesnt exist: %v", err)
			return []model.Todo{}, 0, 0, err
		}
		reminds = append(reminds, remind)
	}

	var nextCursor int

	if len(reminds) > 0 {
		nextCursor = reminds[len(reminds)-1].ID
	} else {
		reminds = []model.Todo{}
	}

	return reminds, totalCount, nextCursor, nil
}

// CreateRemind  store new remind entity to DB PostgresSQL
func (s *TodoStorage) CreateRemind(ctx context.Context, todo model.Todo) (model.Todo, error) {
	var createdTodo model.Todo

	const sql = `INSERT INTO reminder.todo ("Title", "Description",  "User", "CreatedAt", "DeadlineAt", "DeadlineNotify", "NotifyPeriod") 
				 VALUES ($1, $2, $3, $4, $5, $6, $7) returning "ID", "Title", "Description", "User", "CreatedAt", "DeadlineAt", "DeadlineNotify", "NotifyPeriod"`
	row := s.Postgres.QueryRow(ctx, sql, todo.Title, todo.Description, todo.UserID, todo.CreatedAt, todo.DeadlineAt, todo.DeadlineNotify, todo.NotifyPeriod)
	err := row.Scan(
		&createdTodo.ID,
		&createdTodo.Title,
		&createdTodo.Description,
		&createdTodo.UserID,
		&createdTodo.CreatedAt,
		&createdTodo.DeadlineAt,
		&createdTodo.DeadlineNotify,
		&createdTodo.NotifyPeriod,
	)
	if err != nil {
		s.logger.Errorf("Error create remind: %v", err)
		return model.Todo{}, err
	}
	return createdTodo, nil
}

// UpdateRemind update remind, can change Description, Completed and FinishedAt if Completed = true
func (s *TodoStorage) UpdateRemind(ctx context.Context, id int, input model.TodoUpdateInput) (model.Todo, error) {
	const sql = `UPDATE reminder.todo SET "Title" = $1, "Description" = $2, "DeadlineAt"=$3, "FinishedAt" = $4, "Completed" = $5, "DeadlineNotify" = $6, "NotifyPeriod" = $7 WHERE "ID" = $8`

	ct, err := s.Postgres.Exec(ctx, sql, input.Title, input.Description, input.DeadlineAt, input.FinishedAt, input.Completed, input.DeadlineNotify, input.NotifyPeriod, id)
	if err != nil {
		s.logger.Printf("unable to update remind %v", err)
		return model.Todo{}, err
	}

	if ct.RowsAffected() == 0 {
		return model.Todo{}, errors.New("remind not found")
	}

	parseDeadline, err := time.Parse(time.RFC3339, input.DeadlineAt)
	if err != nil {
		return model.Todo{}, err
	}

	var deadlinePeriodNotify []time.Time
	if len(input.NotifyPeriod) > 0 {
		for _, i := range input.NotifyPeriod {
			parseDeadlineNotifyPeriod, err := time.Parse(time.RFC3339, i)
			if err != nil {
				return model.Todo{}, err
			}
			deadlinePeriodNotify = append(deadlinePeriodNotify, parseDeadlineNotifyPeriod)
		}
	}

	var todo model.Todo
	todo.ID = id
	todo.Title = input.Title
	todo.Description = input.Description
	todo.DeadlineAt = parseDeadline
	todo.FinishedAt = input.FinishedAt
	todo.Completed = input.Completed
	todo.DeadlineNotify = input.DeadlineNotify
	todo.NotifyPeriod = deadlinePeriodNotify

	return todo, nil
}

// UpdateNotification update Notificated field
func (s *TodoStorage) UpdateNotification(ctx context.Context, id int, dao model.NotificationDAO) error {
	sql := `UPDATE reminder.todo SET "Notificated" = $1 WHERE "ID" = $2`

	ct, err := s.Postgres.Exec(ctx, sql, dao.Notificated, id)
	if err != nil {
		s.logger.Printf("unable to update notificated status %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("remind not found")
	}

	return nil
}

// UpdateStatus update Completed field
func (s *TodoStorage) UpdateStatus(ctx context.Context, id int, updateInput model.TodoUpdateStatusInput) error {
	const sql = `UPDATE reminder.todo SET "FinishedAt" = $1, "Completed" = $2 WHERE "ID" = $3`

	ct, err := s.Postgres.Exec(ctx, sql, updateInput.FinishedAt, updateInput.Completed, id)
	if err != nil {
		s.logger.Printf("unable to update status %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("remind not found")
	}

	return nil
}

// DeleteRemind deletes remind from DB
func (s *TodoStorage) DeleteRemind(ctx context.Context, id int) error {
	const sql = `DELETE FROM reminder.todo WHERE "ID" = $1`
	res, err := s.Postgres.Exec(ctx, sql, id)
	if err != nil {
		s.logger.Errorf("Error delete remind: %v", err)
		return model.ErrDeleteFailed
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Errorf("error don't found remind: %v", err)
		return model.ErrCantFindRemindWithID
	}

	return nil
}

// GetRemindByID takes out one remind from PostgreSQL by id
func (s *TodoStorage) GetRemindByID(ctx context.Context, id int) (model.Todo, error) {
	var todo model.Todo

	const sql = `SELECT "ID", "Title", "Description", "User", "CreatedAt", "DeadlineAt", "Completed", "FinishedAt", "Notificated" FROM reminder.todo
    WHERE "ID" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, id)

	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.DeadlineAt,
		&todo.Completed,
		&todo.FinishedAt,
		&todo.Notificated,
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

func (s *TodoStorage) GetRemindsForNotification(ctx context.Context) ([]model.NotificationRemind, error) {
	reminds := []model.NotificationRemind{}

	for i := 1; i <= 3; i++ {
		t := time.Now().AddDate(0, 0, i).Format("2006-01-02 15:04:05")
		tn := time.Now().Format("2006-01-02 15:04:05")

		sql := fmt.Sprintf(`SELECT t."ID", t."Description", t."Title", t."DeadlineAt", t."User" from reminder.todo t 
INNER JOIN reminder.users_configs u on u."ID" = t."User" 
WHERE t."DeadlineAt" BETWEEN '%s' AND '%s' 
AND t."Completed" = false 
AND t."Notificated" = false
AND u."Notification" = true
AND u."Period" = %d`, tn, t, i)

		rows, err := s.Postgres.Query(ctx, sql)
		if err != nil {
			s.logger.Errorf("error to select reminds for notification: %v", err)
			return nil, err
		}

		for rows.Next() {
			var remind model.NotificationRemind

			if err := rows.Scan(
				&remind.ID,
				&remind.Description,
				&remind.Title,
				&remind.DeadlineAt,
				&remind.UserID,
			); err != nil {
				s.logger.Errorf("remind doesn't exist: %v", err)
				return nil, err
			}
			reminds = append(reminds, remind)
		}

		rows.Close()
	}

	return reminds, nil
}

func (s *TodoStorage) GetRemindsForDeadlineNotification(ctx context.Context) ([]model.NotificationRemind, string, error) {
	var reminds []model.NotificationRemind
	tn := time.Now().Truncate(time.Minute).Format(time.RFC3339)

	sql := fmt.Sprintf(`SELECT t."ID", t."Description", t."Title", t."DeadlineAt", t."User" from reminder.todo t 
INNER JOIN reminder.users_configs u on u."ID" = t."User" 
WHERE t."NotifyPeriod" @> ARRAY['%s']::TIMESTAMP[] 
AND t."Completed" = false 
AND t."DeadlineNotify" = true`, tn)

	rows, err := s.Postgres.Query(ctx, sql)
	if err != nil {
		s.logger.Errorf("error to select deadline reminds for notification: %v", err)
		return nil, "", err
	}

	defer rows.Close()

	for rows.Next() {
		var remind model.NotificationRemind

		if err := rows.Scan(
			&remind.ID,
			&remind.Description,
			&remind.Title,
			&remind.DeadlineAt,
			&remind.UserID,
		); err != nil {
			s.logger.Errorf("remind doesn't exist: %v", err)
			return nil, "", err
		}
		reminds = append(reminds, remind)
	}

	return reminds, tn, nil
}

func (s *TodoStorage) UpdateNotifyPeriod(ctx context.Context, id int, timeToDelete string) error {
	sql := fmt.Sprintf(`UPDATE reminder.todo SET "NotifyPeriod" = array_remove("NotifyPeriod", '%s')
WHERE "ID" = '%d'`, timeToDelete, id)

	ct, err := s.Postgres.Exec(ctx, sql)
	if err != nil {
		s.logger.Printf("unable to update remind notifier period %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("remind not found")
	}

	return nil
}

package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/red-rocket-software/reminder-go/internal/user/domain"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

var ErrCantGetUserFromDB = errors.New("cannot get user from database")

// UserStorage handles database communication with PostgreSQL.
type UserStorage struct {
	// Postgres database.PGX
	Postgres *pgxpool.Pool
	// Logrus logger
	logger *logging.Logger
}

// NewUserStorage  return new UserStorage with Postgres pool and logger
func NewUserStorage(postgres *pgxpool.Pool, logger *logging.Logger) UserRepo {
	return &UserStorage{Postgres: postgres, logger: logger}
}

func (s *UserStorage) CreateUser(ctx context.Context, user model.User) (int, error) {
	var id int

	const sql = `INSERT INTO users ("Name", "Email", "Password", "Provider", "CreatedAt", "UpdatedAt")
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING "ID"`

	row := s.Postgres.QueryRow(ctx, sql, user.Name, user.Email, user.Password, user.Provider, user.CreatedAt, user.UpdatedAt)

	err := row.Scan(&id)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return 0, err
		}
		s.logger.Errorf("Error create user: %v", err)
		return 0, err
	}
	return id, nil
}

func (s *UserStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	const sql = `SELECT * FROM users WHERE "Email" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Notification,
		&user.Period,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, errors.New("no rows in result set")
	}

	if err != nil {
		s.logger.Errorf("cannot get user from database: %v\n", err)
		return model.User{}, ErrCantGetUserFromDB
	}
	return user, nil
}

func (s *UserStorage) UpdateUser(ctx context.Context, id int, input model.User) error {
	const sql = `UPDATE users SET "Name" = $1, "Email" = $2, "Password" = $3, "Provider" = $4, "CreatedAt" = $5, "UpdatedAt" = $6, "Notification" = $7 WHERE "ID" = $8`

	ct, err := s.Postgres.Exec(ctx, sql, input.Name, input.Email, input.Password, input.Provider, input.CreatedAt, input.UpdatedAt, input.Notification, id)

	if err != nil {
		s.logger.Errorf("unable to update user %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStorage) UpdateUserNotification(ctx context.Context, id int, input model.NotificationUserInput) error {
	var sql string
	if input.Notification != nil {
		sql = fmt.Sprintf(`UPDATE users SET "Notification" = '%t', "Period" = '%d' WHERE "ID" = '%d'`, *input.Notification, input.Period, id)
	} else {
		sql = fmt.Sprintf(`UPDATE users SET "Period" = '%d' WHERE "ID" = '%d'`, input.Period, id)
	}

	ct, err := s.Postgres.Exec(ctx, sql)

	if err != nil {
		s.logger.Errorf("unable to update user %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStorage) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User

	const sql = `SELECT "ID", "Name", "Email", "Password", "Provider", "CreatedAt", "UpdatedAt" FROM users WHERE "ID" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, errors.New("no rows in result set")
	}

	if err != nil {
		s.logger.Printf("cannot get user from database: %v\n", err)
		return model.User{}, ErrCantGetUserFromDB
	}
	return user, nil
}

// DeleteUser deletes user from DB
func (s *UserStorage) DeleteUser(ctx context.Context, id int) error {
	const sql = `DELETE FROM users WHERE "ID" = $1`
	res, err := s.Postgres.Exec(ctx, sql, id)

	if err != nil {
		s.logger.Errorf("error don't found user: %v", err)
		return ErrCantGetUserFromDB
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Errorf("Error delete user: %v", err)
		return ErrCantGetUserFromDB
	}

	return nil
}

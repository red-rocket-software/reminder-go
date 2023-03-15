package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
)

func (s *TodoStorage) CreateUser(ctx context.Context, user model.User) (int, error) {
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

func (s *TodoStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	const sql = `SELECT "ID", "Name", "Email", "Password", "Provider", "CreatedAt", "UpdatedAt", "Period" FROM users WHERE "Email" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
		&user.CreatedAt,
		&user.UpdatedAt,
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

func (s *TodoStorage) UpdateUser(ctx context.Context, id int, input model.User) error {
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

func (s *TodoStorage) UpdateUserNotification(ctx context.Context, id int, status bool, period int) error {
	var sql = fmt.Sprintf(`UPDATE users SET "Notification" = '%t', "Period" = '%d' WHERE "ID" = '%d'`, status, period, id)

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

func (s *TodoStorage) GetUserByID(ctx context.Context, id int) (model.User, error) {
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
func (s *TodoStorage) DeleteUser(ctx context.Context, id int) error {
	const sql = `DELETE FROM users WHERE "ID" = $1`
	res, err := s.Postgres.Exec(ctx, sql, id)

	if err != nil {
		s.logger.Errorf("error don't found user: %v", err)
		return ErrCantFindRemindWithID
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Errorf("Error delete user: %v", err)
		return ErrDeleteFailed
	}

	return nil
}

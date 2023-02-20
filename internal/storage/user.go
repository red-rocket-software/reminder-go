package storage

import (
	"context"
	"errors"
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

	const sql = `SELECT "ID", "Name", "Email", "Password", "Provider", "CreatedAt", "UpdatedAt" FROM users WHERE "Email" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, email)

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
		s.logger.Errorf("cannot get user from database: %v\n", err)
		return model.User{}, ErrCantGetUserFromDB
	}
	return user, nil
}

func (s *TodoStorage) UpdateUser(ctx context.Context, id int, input model.User) error {
	const sql = `UPDATE users SET "Name" = $1, "Email" = $2, "Password" = $3, "Provider" = $4, "CreatedAt" = $5, "UpdatedAt" = $6`

	ct, err := s.Postgres.Exec(ctx, sql, input.Name, input.Email, input.Password, input.Provider, input.CreatedAt, input.UpdatedAt)

	if err != nil {
		s.logger.Errorf("unable to update user %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *TodoStorage) GetUserById(ctx context.Context, id int) (model.User, error) {
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

package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/red-rocket-software/reminder-go/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// NewClient Create Postgres pgx connection with attempts
func NewClient(ctx context.Context, maxAttemps int, cfg config.Config) (con *pgx.Conn, err error) {

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		con, err = pgx.Connect(ctx, dsn)
		if err != nil {
			fmt.Println("failed to connect to postgesql")
			return err
		}
		return nil
	}, maxAttemps, 5*time.Second)
	if err != nil {
		log.Fatalln("error do with tries postgresql")
	}
	return con, nil
}

// DoWithTries  provide attempts to connect db
func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}
		return nil
	}
	return
}

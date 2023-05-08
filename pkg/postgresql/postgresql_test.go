package postgresql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/red-rocket-software/reminder-go/config"
)

func TestNewPGXPoolErrors(t *testing.T) {
	t.Parallel()
	cfg := config.Config{}

	cfg.Postgres.Host = "wrong"
	cfg.Postgres.Username = "wrong"
	cfg.Postgres.Port = "1234"
	cfg.Postgres.Database = "wrong"

	type args struct {
		ctx      context.Context
		attempts int
		config   config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *pgxpool.Pool
		wantErr bool
	}{
		{
			name: "invalid_connection_string",
			args: args{
				ctx:      context.Background(),
				attempts: 5,
				config:   cfg,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.ctx, tt.args.attempts, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPGXPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("NewPGXPool() = %v, want nil", got)
			}
		})
	}
}

func TestDoWithAttempts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		attempts := 0
		err := DoWithTries(func() error {
			attempts++
			return nil
		}, 3, time.Millisecond)
		if err != nil {
			t.Errorf("Expected nil error, got: %v", err)
		}
		if attempts != 1 {
			t.Errorf("Expected 1 attempt, got: %v", attempts)
		}
	})

	t.Run("All attempts fail", func(t *testing.T) {
		attempts := 0
		err := DoWithTries(func() error {
			attempts++
			return fmt.Errorf("failed attempt")
		}, 3, time.Millisecond)
		if err == nil {
			t.Errorf("Expected non-nil error, got nil")
		}
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got: %v", attempts)
		}
	})

	t.Run("Succeeds on second attempt", func(t *testing.T) {
		attempts := 0
		err := DoWithTries(func() error {
			attempts++
			if attempts == 2 {
				return nil
			}
			return fmt.Errorf("failed attempt")
		}, 3, time.Millisecond)
		if err != nil {
			t.Errorf("Expected nil error, got: %v", err)
		}
		if attempts != 2 {
			t.Errorf("Expected 2 attempts, got: %v", attempts)
		}
	})
}

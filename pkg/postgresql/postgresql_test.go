package postgresql

import (
	"context"
	"testing"

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

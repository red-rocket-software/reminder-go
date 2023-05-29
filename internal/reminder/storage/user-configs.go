package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

var _ model.ConfigRepository = (*ConfigsStorage)(nil)

// ConfigsStorage handles database communication with PostgreSQL.
type ConfigsStorage struct {
	// Postgres database.PGX
	Postgres *pgxpool.Pool
	// Logrus logger
	logger *logging.Logger
}

// NewConfigsStorage  return new ConfigsStorage with Postgres pool and logger
func NewConfigsStorage(postgres *pgxpool.Pool, logger *logging.Logger) model.ConfigRepository {
	return &ConfigsStorage{Postgres: postgres, logger: logger}
}

// UpdateUserConfig update user_configs. Changes notification or period
func (s *ConfigsStorage) UpdateUserConfig(ctx context.Context, id string, input model.UserConfigs) error {
	tn := time.Now()
	const sql = `UPDATE reminder.users_configs SET "Notification" = $1, "Period" = $2, "UpdatedAt" = $3 WHERE "ID" = $4`

	ct, err := s.Postgres.Exec(ctx, sql, input.Notification, input.Period, tn, id)

	if err != nil {
		s.logger.Errorf("unable to update user-config %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("user configs not found")
	}

	return nil
}

// GetUserConfigs returns user configs from database
func (s *ConfigsStorage) GetUserConfigs(ctx context.Context, userID string) (model.UserConfigs, error) {
	var configs model.UserConfigs

	const sql = `SELECT "ID", "Notification", "Period", "CreatedAt", "UpdatedAt"  FROM reminder.users_configs
    WHERE "ID" = $1 LIMIT 1`

	row := s.Postgres.QueryRow(ctx, sql, userID)

	err := row.Scan(
		&configs.ID,
		&configs.Notification,
		&configs.Period,
		&configs.CreatedAt,
		&configs.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.UserConfigs{}, nil
	}
	if err != nil {
		s.logger.Printf("cannot get user-configs from database: %v\n", err)
		return model.UserConfigs{}, errors.New("cannot get user-configs from database")
	}

	return configs, nil
}

// CreateUserConfigs  store new remind entity to DB PostgresSQL
func (s *ConfigsStorage) CreateUserConfigs(ctx context.Context, userID string) (model.UserConfigs, error) {
	var userConfig model.UserConfigs

	userConfig.ID = userID
	userConfig.Notification = false
	userConfig.Period = 2
	userConfig.CreatedAt = time.Now()

	const sql = `INSERT INTO reminder.users_configs ("ID", "Notification",  "Period", "CreatedAt") 
				 VALUES ($1, $2, $3, $4) returning "ID", "Notification",  "Period", "CreatedAt", "UpdatedAt"`
	row := s.Postgres.QueryRow(ctx, sql, userConfig.ID, userConfig.Notification, userConfig.Period, userConfig.CreatedAt)
	err := row.Scan(
		&userConfig.ID,
		&userConfig.Notification,
		&userConfig.Period,
		&userConfig.CreatedAt,
		&userConfig.UpdatedAt,
	)
	log.Print("CreatedAt ", userConfig.CreatedAt)
	if err != nil {
		s.logger.Errorf("Error create userConfigs: %v", err)
		return model.UserConfigs{}, err
	}
	return userConfig, nil
}

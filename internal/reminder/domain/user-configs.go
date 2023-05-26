package domain

import (
	"context"
	"time"
)

type UserConfigs struct {
	ID           string     `json:"ID"`
	Notification bool       `json:"notification"`
	Period       int        `json:"period"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

//go:generate mockgen -source=user-configs.go -destination=mocks/configsStorage.go

type ConfigRepository interface {
	GetUserConfigs(ctx context.Context, userID string) (UserConfigs, error)
	CreateUserConfigs(ctx context.Context, userID string) (UserConfigs, error)
	UpdateUserConfig(ctx context.Context, id string, input UserConfigs) error
}

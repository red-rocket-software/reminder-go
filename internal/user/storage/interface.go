package storage

import (
	"context"

	model "github.com/red-rocket-software/reminder-go/internal/user/domain"
)

//go:generate mockgen -source=interface.go -destination=mocks/userStorage.go

type UserRepo interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, id int, input model.User) error
	GetUserByID(ctx context.Context, id int) (model.User, error)
	UpdateUserNotification(ctx context.Context, id int, input model.NotificationUserInput) error
	DeleteUser(ctx context.Context, id int) error
}

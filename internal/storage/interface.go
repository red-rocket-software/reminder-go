package storage

import (
	"context"
	"errors"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/pkg/pagination"
)

var (
	ErrDeleteFailed         = errors.New("delete failed")
	ErrCantFindRemind       = errors.New("cannot get product from database")
	ErrCantFindRemindWithID = errors.New("cannot find remind with such id")
	ErrCantGetUserFromDB    = errors.New("cannot get user from database")
)

//go:generate mockgen -source=interface.go -destination=mocks/storage.go

type ReminderRepo interface {
	GetAllReminds(ctx context.Context, params pagination.Page, userID int) ([]model.Todo, int, int, error)
	CreateRemind(ctx context.Context, todo model.Todo) (model.Todo, error)
	UpdateRemind(ctx context.Context, id int, input model.TodoUpdateInput) (model.Todo, error)
	UpdateStatus(ctx context.Context, id int, updateInput model.TodoUpdateStatusInput) error
	UpdateNotification(ctx context.Context, id int, dao model.NotificationDAO) error
	DeleteRemind(ctx context.Context, id int) error
	GetRemindByID(ctx context.Context, id int) (model.Todo, error)
	GetCompletedReminds(ctx context.Context, params Params, userID int) ([]model.Todo, int, int, error)
	GetNewReminds(ctx context.Context, params pagination.Page, userID int) ([]model.Todo, int, int, error)
	Truncate() error
	SeedTodos() ([]model.Todo, error)
	SeedUser() (int, error)
	GetRemindsForNotification(ctx context.Context) ([]model.NotificationRemind, error)
	GetRemindsForDeadlineNotification(ctx context.Context) ([]model.NotificationRemind, string, error)
	UpdateNotifyPeriod(ctx context.Context, id int, timeToDelete string) error

	// user methods
	CreateUser(ctx context.Context, input model.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, id int, input model.User) error
	GetUserByID(ctx context.Context, id int) (model.User, error)
	UpdateUserNotification(ctx context.Context, id int, input model.NotificationUserInput) error
	DeleteUser(ctx context.Context, id int) error
}

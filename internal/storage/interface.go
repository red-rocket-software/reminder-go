package storage

import (
	"context"
	"errors"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
)

var (
	ErrDeleteFailed   = errors.New("delete failed")
	ErrCantFindRemind = errors.New("cannot get product from database")
)

//go:generate mockgen -source=interface.go -destination=mocks/storage.go

type ReminderRepo interface {
	GetAllReminds(ctx context.Context, fetchParams FetchParam) ([]model.Todo, int, error)
	CreateRemind(ctx context.Context, todo model.Todo) (int, error)
	UpdateRemind(ctx context.Context, id int, input model.TodoUpdate) error
	DeleteRemind(ctx context.Context, id int) error
	GetRemindByID(ctx context.Context, id int) (model.Todo, error)
	GetComplitedReminds(ctx context.Context, params FetchParam) ([]model.Todo, int, error)
	GetNewReminds(ctx context.Context, params FetchParam) ([]model.Todo, int, error)
	Truncate() error
	SeedTodos() ([]model.Todo, error)
}

package storage

import (
	"context"
	"errors"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
)

var (
	ErrDeleteFailed   = errors.New("delete failed")
	ErrCantFindRemind = errors.New("can't find remind")
)

type ReminderRepo interface {
	GetAllReminds(ctx context.Context) ([]model.Todo, error)
	CreateRemind(ctx context.Context, todo model.Todo) error
	UpdateRemind(ctx context.Context, id int) (model.Todo, error)
	DeleteRemind(ctx context.Context, id string) error
	GetRemindByID(ctx context.Context, id int) (model.Todo, error)
	GetComplitedReminds(ctx context.Context) ([]model.Todo, error)
	GetNewReminds(ctx context.Context) ([]model.Todo, error)
}

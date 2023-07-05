package domain

import (
	"context"
	"errors"
	"time"

	"github.com/red-rocket-software/reminder-go/pkg/utils"
)

var (
	ErrDeleteFailed         = errors.New("error delete remind")
	ErrCantFindRemindWithID = errors.New("can't find remind")
)

type Todo struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	UserID         string      `json:"user_id"`
	CreatedAt      time.Time   `json:"created_at"`
	DeadlineAt     time.Time   `json:"deadline_at"`
	FinishedAt     *time.Time  `json:"finished_at,omitempty"`
	Completed      bool        `json:"completed"`
	Notificated    bool        `json:"notificated"`
	DeadlineNotify *bool       `json:"deadline_notify"`
	NotifyPeriod   []time.Time `json:"notify_period"`
}

type TodoInput struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	DeadlineAt     string   `json:"deadline_at"`
	CreatedAt      string   `json:"created_at"`
	DeadlineNotify *bool    `json:"deadline_notify"`
	NotifyPeriod   []string `json:"notify_period"`
}

type TodoUpdateInput struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
	Completed      bool       `json:"completed"`
	Notificated    bool       `json:"notificated"`
	DeadlineAt     string     `json:"deadline_at"`
	DeadlineNotify *bool      `json:"deadline_notify"`
	NotifyPeriod   []string   `json:"notify_period"`
}

type TodoResponse struct {
	Todos    []Todo         `json:"todos"`
	Count    int            `json:"count"`
	PageInfo utils.PageInfo `json:"pageInfo"`
}

type TodoUpdateStatusInput struct {
	Completed  bool       `json:"completed"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type NotificationRemind struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DeadlineAt  time.Time `json:"deadline_at"`
	UserID      string    `json:"user_id"`
	MessageType string    `json:"message_type"`
}

type NotificationDAO struct {
	Notificated bool
}

type TimeRangeFilter struct {
	StartRange string
	EndRange   string
}

type FetchParams struct {
	utils.Page
	TimeRangeFilter
	FilterByDate  string //createdAt or deadlineAt
	FilterBySort  string // ASC or DESC
	FilterByQuery string // current, all or completed

}

//go:generate mockgen -source=todo.go -destination=mocks/todoStorage.go
type TodoRepository interface {
	GetReminds(ctx context.Context, params FetchParams, userID string) ([]Todo, int, int, error)
	CreateRemind(ctx context.Context, todo Todo) (Todo, error)
	UpdateRemind(ctx context.Context, id int, input TodoUpdateInput) (Todo, error)
	UpdateStatus(ctx context.Context, id int, updateInput TodoUpdateStatusInput) error
	UpdateNotification(ctx context.Context, id int, dao NotificationDAO) error
	DeleteRemind(ctx context.Context, id int) error
	GetRemindByID(ctx context.Context, id int) (Todo, error)
	UpdateNotifyPeriod(ctx context.Context, id int, timeToDelete string) error
	GetRemindsForNotification(ctx context.Context) ([]NotificationRemind, error)
	GetRemindsForDeadlineNotification(ctx context.Context) ([]NotificationRemind, string, error)
}

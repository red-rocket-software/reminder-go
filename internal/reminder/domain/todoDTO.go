package domain

import (
	"time"

	"github.com/red-rocket-software/reminder-go/pkg/pagination"
)

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
	Todos    []Todo              `json:"todos"`
	Count    int                 `json:"count"`
	PageInfo pagination.PageInfo `json:"pageInfo"`
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
	UserID      int       `json:"user_id"`
}

type NotificationDAO struct {
	Notificated bool
}

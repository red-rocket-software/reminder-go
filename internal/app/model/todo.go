package model

import (
	"time"

	"github.com/red-rocket-software/reminder-go/pkg/pagination"
)

type Todo struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	UserID      int        `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	DeadlineAt  time.Time  `json:"deadline_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	Completed   bool       `json:"completed"`
	Notificated bool       `json:"notificated"`
}

type TodoInput struct {
	Description string `json:"description"`
	DeadlineAt  string `json:"deadline_at"`
	CreatedAt   string `json:"created_at"`
}

type TodoUpdateInput struct {
	Description string     `json:"description"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	Completed   bool       `json:"completed"`
	Notificated bool       `json:"notificated"`
	DeadlineAt  string     `json:"deadline_at"`
}

type TodoResponse struct {
	Todos    []Todo              `json:"todos"`
	PageInfo pagination.PageInfo `json:"pageInfo"`
}

type TodoUpdateStatusInput struct {
	Completed  bool       `json:"completed"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type NotificationRemind struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	DeadlineAt  time.Time `json:"deadline_at"`
	UserID      int       `json:"user_id"`
}

type NotificationDAO struct {
	Notificated bool
}

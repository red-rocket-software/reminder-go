package model

import (
	"time"

	"github.com/red-rocket-software/reminder-go/pkg/pagination"
)

type Todo struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	DeadlineAt  time.Time  `json:"deadline_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	Completed   bool       `json:"completed"`
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

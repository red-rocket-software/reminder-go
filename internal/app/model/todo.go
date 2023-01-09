package model

import (
	"time"
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
}

type TodoUpdate struct {
	Description string     `json:"description"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	Completed   bool       `json:"completed"`
}
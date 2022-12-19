package model

import (
	"time"
)

type Todo struct {
	ID          string    `json:"_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	DeadlineAt  time.Time `json:"deadlineAt"`
	FinishedAt  time.Time `json:"finishedAt,omitempty"`
	Completed   bool      `json:"completed"`
}

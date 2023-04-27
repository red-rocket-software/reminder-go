package domain

import (
	"time"
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

type UserConfigs struct {
	ID           string     `json:"ID"`
	Notification bool       `json:"notification"`
	Period       int        `json:"period"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

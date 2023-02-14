package model

import (
	"time"
)

type User struct {
	ID       int    `json: "id"`
	Name     string `json: "name"`
	Email    string `json: "email"`
	Password string `json: "password"`

	Provider  string    `json: "provider`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json: "updated_at, omitempty"`
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

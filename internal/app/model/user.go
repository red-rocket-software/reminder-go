package model

import (
	"time"
)

type user struct {
	ID                int       `json: "id"`
	Nickname          string    `json: "nickname"`
	Email             string    `json: "email"`
	EncryptedPassword string    `json: "encrypted_password"`
	UpdatedAt         string    `json: "updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}

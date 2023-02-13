package model

import (
	"time"
)

type user struct {
	ID                int       `json: "id"`
	Nickname          string    `json: "nickname"`
	Email             string    `json: "email"`
	CreatedAt         time.Time `json:"created_at"`
	EncryptedPassword string    `json: "encrypted_password"`
}

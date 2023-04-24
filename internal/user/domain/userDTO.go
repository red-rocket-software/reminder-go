package domain

import "time"

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ToResponseUser(user User) UserResponse {
	userResponse := UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Period:       user.Period,
		Notification: user.Notification,
	}

	return userResponse
}

type UserResponse struct {
	ID           int       `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Notification bool      `json:"notification"`
	Period       *int      `json:"period"`
}

type NotificationUserInput struct {
	Notification *bool `json:"notification,omitempty"`
	Period       int   `json:"period"`
}

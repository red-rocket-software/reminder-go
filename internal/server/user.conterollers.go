package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/internal/server/auth"
	"github.com/red-rocket-software/reminder-go/utils"
)

func (server *Server) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload model.RegisterUserInput

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	now := time.Now()

	newUser := model.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  payload.Password,
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = model.Validate(newUser, "")
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err := server.TodoStorage.SaveUser(server.ctx, newUser)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusCreated, "User is successfully created")
}

func (server *Server) SignInUser(w http.ResponseWriter, r *http.Request) {
	var payload model.LoginUserInput

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user model.User
	result, err := server.TodoStorage.GetUserById(server.ctx, payload.Email)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	if user.Provider == "Google" {
		utils.JSONFormat(w, http.StatusUnauthorized, fmt.Sprintf("Use %v OAuth instead", user.Provider))
		return
	}

	token, err := auth.GenerateToken(server.config.Auth.TokenExpiredIn, user.ID, server.config.Auth.JwtSecret)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = server.config.Auth.TokenMaxAge
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

	utils.JSONFormat(w, http.StatusCreated, "Successful logIn")
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{}
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = -1
	cookie.Secure = false
	cookie.HttpOnly = true

	utils.JSONFormat(w, http.StatusOK, "Success")
}

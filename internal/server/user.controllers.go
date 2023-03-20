package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/utils"
)

// GetMe godoc
//
//	@Description	GetMe
//	@Summary		fetch current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	model.User
//
//
//	@Security		BasicAuth
//	@Router			/fetchMe [get]
func (server *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value("currentUser").(model.User)

	utils.JSONFormat(w, http.StatusOK, currentUser)
}

// UpdateUserNotification godoc
//
//	@Description	UpdateUserNotification
//	@Summary		update user notification status
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			input	body		model.NotificationUserInput	true	"update info"
//	@Success		200		{string}	string						"user notification status successfully updated"
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/status/{id} [put]
func (server *Server) UpdateUserNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	var input model.NotificationUserInput

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if input.Period == 0 {
		input.Period = 2
	}

	err = server.TodoStorage.UpdateUserNotification(server.ctx, uID, input)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, "user notification status successfully updated")
}

// DeleteUser godoc
//
//	@Description	DeleteUser
//	@Summary		delete user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"id"
//	@Success		204	{string}	string	"user with id:1 successfully deleted"
//
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/user{id} [delete]
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	// deleting remind from db
	if err := server.TodoStorage.DeleteUser(server.ctx, userID); err != nil {
		if errors.Is(err, storage.ErrDeleteFailed) {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}
		if errors.Is(err, storage.ErrCantFindRemindWithID) {
			utils.JSONError(w, http.StatusNotFound, err)
		}
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = -1
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

	successMsg := fmt.Sprintf("user with id:%d successfully deleted", userID)

	w.Header().Set("Content-Type", "application/json")
	utils.JSONFormat(w, http.StatusNoContent, successMsg)
}

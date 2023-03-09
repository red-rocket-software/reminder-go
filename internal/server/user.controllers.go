package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/utils"
)

func (server *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value("currentUser").(model.User)

	utils.JSONFormat(w, http.StatusOK, currentUser)
}

func (server *Server) UpdateUserNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	type userNotification struct {
		notification bool
	}

	var input userNotification

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.TodoStorage.UpdateUserNotification(server.ctx, uID, input.notification)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, "user notification status successfully updated")
}

package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/utils"
)

type TodoHandlers interface {
	GetAllReminds(w http.ResponseWriter, r *http.Request)
	GetRemindById(w http.ResponseWriter, r *http.Request)
	AddRemind(w http.ResponseWriter, r *http.Request)
	UpdateRemind(w http.ResponseWriter, r *http.Request)
	DeleteRemind(w http.ResponseWriter, r *http.Request)
	GetComplitedReminds(w http.ResponseWriter, r *http.Request)
	GetCurrentReminds(w http.ResponseWriter, r *http.Request)
}

func (s *Server) AddRemind(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input model.TodoInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if input.Description == "" {
		utils.JsonError(w, http.StatusUnprocessableEntity, errors.New("nothing to save"))
		return
	}

	var todo model.Todo

	parseTime, err := time.Parse("2006-01-02", input.DeadlineAt)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	todo.CreatedAt = time.Now()
	todo.Description = input.Description
	todo.DeadlineAt = parseTime

	err = s.TodoStorage.CreateRemind(s.ctx, todo)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JsonFormat(w, http.StatusCreated, "remind successfully created")
}

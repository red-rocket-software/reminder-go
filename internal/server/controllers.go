package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/internal/storage"
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

// AddRemind gets remind from user input, decode and sent to DB. Simple validation - no empty field Description.
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

func (s *Server) DeleteRemind(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	remindID := vars["id"]

	// Check if the remind exist
	//TODO you should imolpement GetRemindByID method in storage
	// _, err := s.TodoStorage.GetRemindByID(s.ctx, remindID)
	// if errors.Is(err, storage.ErrCantFindRemind) {
	// 	utils.JsonError(w, http.StatusInternalServerError, err)
	// 	return
	// } else {
	// 	utils.JsonError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	// deleting remind from db
	if err := s.TodoStorage.DeleteRemind(s.ctx, remindID); err != nil {
		if errors.Is(err, storage.ErrDeleteFailed) {
			utils.JsonError(w, http.StatusInternalServerError, err)
			return
		} else {
			utils.JsonError(w, http.StatusInternalServerError, err)
			return
		}
	}

	successMsg := fmt.Sprintf("remind with id:%s successfully deleted", remindID)

	utils.JsonFormat(w, http.StatusCreated, successMsg)
}

func (s *Server) GetRemindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	rId, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	todo, err := s.TodoStorage.GetRemindByID(s.ctx, rId)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}
	utils.JsonFormat(w, http.StatusOK, todo)
}

//UpdateRemind update Description field and Completed if true changes FinishedAt on time.Now
func (s *Server) UpdateRemind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	rId, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	var input model.TodoUpdate

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if input.Completed == true {
		input.FinishedAt = time.Now()
	}

	if input.Description == "" {
		utils.JsonError(w, http.StatusUnprocessableEntity, errors.New("description is empty"))
		return
	}

	err = s.TodoStorage.UpdateRemind(s.ctx, rId, input)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JsonFormat(w, http.StatusOK, "remind successfully updated")
}

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

// GetAllReminds makes request to DB for all reminds. Works with cursor pagination
func (s *Server) GetAllReminds(w http.ResponseWriter, r *http.Request) {
	// scan for limit in parameters
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		utils.JsonError(w, http.StatusBadRequest, errors.New("limit parameter is invalid"))
		return
	}
	if limit == 0 {
		limit = 10
	}

	// scan for cursor in parameters
	cursorStr := r.URL.Query().Get("cursor")
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil && cursorStr != "" {
		utils.JsonError(w, http.StatusBadRequest, errors.New("cursor parameter is invalid"))
		return
	}

	//inititalize fetchParameters
	fetchParams := storage.FetchParams{
		Limit:  uint64(limit),
		Cursor: uint64(cursor),
	}

	reminds, nextCursor, err := s.TodoStorage.GetAllReminds(s.ctx, fetchParams)
	if err != nil && cursorStr != "" {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-NextCursor", fmt.Sprintf("%d", nextCursor))

	utils.JsonFormat(w, http.StatusOK, reminds)
}

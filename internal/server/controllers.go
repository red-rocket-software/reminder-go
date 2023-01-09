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

	dParseTime, err := time.Parse("2006-01-02", input.DeadlineAt)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	todo.CreatedAt = time.Now()
	todo.Description = input.Description
	todo.DeadlineAt = dParseTime

	err = s.TodoStorage.CreateRemind(s.ctx, todo)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JsonFormat(w, http.StatusCreated, "remind successfully created")
}

func (s *Server) DeleteRemind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	remindID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	// Check if the remind exist
	_, err = s.TodoStorage.GetRemindByID(s.ctx, remindID)
	if errors.Is(err, storage.ErrCantFindRemind) {
		utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}

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

	successMsg := fmt.Sprintf("remind with id:%d successfully deleted", remindID)

	w.Header().Set("Content-Type", "application/json")
	utils.JsonFormat(w, http.StatusCreated, successMsg)
}

// GetCurrentReminds handle get current reminds. First url should be like: http://localhost:8000/current?limit=5
// the next we should write cursor from prev. headers X-Nextcursor:  http://localhost:8000/current?limit=5&cursor=33
func (s *Server) GetCurrentReminds(w http.ResponseWriter, r *http.Request) {
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if err != nil && strLimit != "" {
		utils.JsonError(w, http.StatusBadRequest, errors.New("limit parameter is invalid, should be positive integer"))
		return
	}

	//if no write limit it will be 5
	if limit == 0 {
		limit = 5
	}

	strCursor := r.URL.Query().Get("cursor")
	cursor, err := strconv.Atoi(strCursor)
	if err != nil && strCursor != "" {
  		utils.JsonError(w, http.StatusBadRequest, errors.New("cursor parameter is invalid"))
		return
	}

	fetchParam := storage.FetchParam{
		Limit:    limit,
		CursorID: cursor,
	}

	reminds, nextCursor, err := s.TodoStorage.GetNewReminds(s.ctx, fetchParam)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}
  
	w.Header().Add("X-NextCursor", fmt.Sprintf("%d", nextCursor))

	utils.JsonFormat(w, http.StatusOK, reminds)
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

// UpdateRemind update Description field and Completed if true changes FinishedAt on time.Now
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

	tn := time.Now()

	if input.Completed == true {
		input.FinishedAt = &tn
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

func (s *Server) GetComplitedReminds(w http.ResponseWriter, r *http.Request) {
	// scan for limit in parameters
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		utils.JsonError(w, http.StatusBadRequest, errors.New("limit parameter is invalid"))
		return
	}
	if limit == 0 {
  		limit = 5
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
		Limit:  limit,
		CursorID: cursor,
	}

	reminds, nextCursor, err := s.TodoStorage.GetComplitedReminds(s.ctx, fetchParams)

	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
  w.Header().Add("X-NextCursor", fmt.Sprintf("%d", nextCursor))

	utils.JsonFormat(w, http.StatusOK, reminds)
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
		Limit:  limit,
		CursorID: cursor,
	}
  
  	reminds, nextCursor, err := s.TodoStorage.GetAllReminds(s.ctx, fetchParams)

	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
  w.Header().Add("X-NextCursor", fmt.Sprintf("%d", nextCursor))

	utils.JsonFormat(w, http.StatusOK, reminds)
}
package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/pkg/pagination"
	"github.com/red-rocket-software/reminder-go/utils"
)

type RemindHandlers interface {
	GetAllReminds(w http.ResponseWriter, r *http.Request)
	GetRemindByID(w http.ResponseWriter, r *http.Request)
	AddRemind(w http.ResponseWriter, r *http.Request)
	UpdateRemind(w http.ResponseWriter, r *http.Request)
	UpdateCompleteStatus(w http.ResponseWriter, r *http.Request)
	DeleteRemind(w http.ResponseWriter, r *http.Request)
	GetCompletedReminds(w http.ResponseWriter, r *http.Request)
	GetCurrentReminds(w http.ResponseWriter, r *http.Request)
}

// AddRemind godoc
//
//	@Description	AddRemind
//	@Summary		create a new remind
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			input	body		model.TodoInput	true	"remind info"
//	@Success		201		{string}	string			"Remind is successfully created"
//
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/remind [post]
func (server *Server) AddRemind(w http.ResponseWriter, r *http.Request) {
	var input model.TodoInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if input.Description == "" || input.DeadlineAt == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("nothing to save"))
		return
	}

	user := r.Context().Value("currentUser").(model.User)

	var todo model.Todo

	deadlineParseTime, err := time.Parse("2006-01-02T15:04", input.DeadlineAt)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	createParseTime, err := time.Parse("02.01.2006, 15:04:05", input.CreatedAt)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	todo.CreatedAt = createParseTime
	todo.Description = input.Description
	todo.DeadlineAt = deadlineParseTime
	todo.UserID = user.ID

	_, err = server.TodoStorage.CreateRemind(server.ctx, todo)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusCreated, "Remind is successfully created")
}

// DeleteRemind godoc
//
//	@Description	DeleteRemind
//	@Summary		delete remind
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"id"
//	@Success		204	{string}	string	"remind with id:1 successfully deleted"
//
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/remind{id} [delete]
func (server *Server) DeleteRemind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	remindID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	// deleting remind from db
	if err := server.TodoStorage.DeleteRemind(server.ctx, remindID); err != nil {
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

	successMsg := fmt.Sprintf("remind with id:%d successfully deleted", remindID)

	w.Header().Set("Content-Type", "application/json")
	utils.JSONFormat(w, http.StatusNoContent, successMsg)
}

// GetCurrentReminds handle get current reminds. First url should be like: http://localhost:8000/current?limit=5
// the next we should write cursor from prev.   http://localhost:8000/current?limit=5&cursor=33
//
//	@Description	GetCurrentReminds
//	@Summary		return a list of current reminds
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		string	true	"limit"
//	@Param			cursor	query		string	true	"cursor"
//	@Success		200		{object}	model.TodoResponse
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/current [get]
func (server *Server) GetCurrentReminds(w http.ResponseWriter, r *http.Request) {
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if err != nil && strLimit != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("limit parameter is invalid, should be positive integer"))
		return
	}

	//if no write limit it will be 5
	if limit == 0 {
		limit = 5
	}

	strCursor := r.URL.Query().Get("cursor")
	cursor, err := strconv.Atoi(strCursor)
	if err != nil && strCursor != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("cursor parameter is invalid"))
		return
	}

	fetchParam := pagination.Page{
		Limit:  limit,
		Cursor: cursor,
	}

	user := r.Context().Value("currentUser").(model.User)

	reminds, nextCursor, err := server.TodoStorage.GetNewReminds(server.ctx, fetchParam, user.ID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	res := model.TodoResponse{
		Todos: reminds,
		PageInfo: pagination.PageInfo{
			Page:       fetchParam,
			NextCursor: nextCursor,
		},
	}

	utils.JSONFormat(w, http.StatusOK, res)
}

// GetRemindByID godoc
//
//	@Description	GetRemindByID
//	@Summary		return a remind by id
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	model.Todo
//
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		404	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/remind/{id} [get]
func (server *Server) GetRemindByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	rID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	todo, err := server.TodoStorage.GetRemindByID(server.ctx, rID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSONError(w, http.StatusNotFound, err)
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSONFormat(w, http.StatusOK, todo)
}

// UpdateRemind update Description field and Completed if true changes FinishedAt on time.Now
//
//	@Description	UpdateRemind
//	@Summary		update remind with given fields
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"id"
//	@Param			input	body		model.TodoUpdateInput	true	"update info"
//	@Success		200		{string}	string					"remind successfully updated"
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/remind/{id} [put]
func (server *Server) UpdateRemind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	rID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	var input model.TodoUpdateInput

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	tn := time.Now()

	if input.Completed {
		input.FinishedAt = &tn
	} else {
		input.FinishedAt = nil
	}

	if input.Description == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("description is empty"))
		return
	}

	err = server.TodoStorage.UpdateRemind(server.ctx, rID, input)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, "remind successfully updated")
}

// UpdateCompleteStatus update Completed field to true
//
//	@Description	UpdateCompleteStatus
//	@Summary		update remind's field "completed"
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			input	body		model.TodoUpdateStatusInput	true	"update info"
//	@Success		200		{string}	string						"remind status updated"
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/status/{id} [put]
func (server *Server) UpdateCompleteStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	rID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	var updateInput model.TodoUpdateStatusInput

	err = json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	tn := time.Now()

	if updateInput.Completed {
		updateInput.FinishedAt = &tn
	}

	err = server.TodoStorage.UpdateStatus(server.ctx, rID, updateInput)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, "remind status updated")
}

// GetCompletedReminds handle get completed reminds.
//
//	@Description	GetCompletedReminds
//	@Summary		return a list of completed reminds
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		string	true	"limit"
//	@Param			cursor	query		string	true	"cursor"
//	@Param			start	query		string	true	"start of time range"
//	@Param			end		query		string	true	"end of time range"
//	@Success		200		{object}	model.TodoResponse
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/completed [get]
func (server *Server) GetCompletedReminds(w http.ResponseWriter, r *http.Request) {
	// scan for limit in parameters
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("limit parameter is invalid, should be positive integer"))
		return
	}
	if limit == 0 {
		limit = 5
	}

	// scan for cursor in parameters
	cursorStr := r.URL.Query().Get("cursor")
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil && cursorStr != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("cursor parameter is invalid"))
		return
	}

	// scan for timeRangeValues in parameters
	rangeStart := r.URL.Query().Get("start")
	rangeEnd := r.URL.Query().Get("end")

	//initialize fetchParameters
	fetchParams := storage.Params{
		Page: pagination.Page{
			Cursor: cursor,
			Limit:  limit,
		},
		TimeRangeFilter: storage.TimeRangeFilter{
			StartRange: rangeStart,
			EndRange:   rangeEnd,
		},
	}

	user := r.Context().Value("currentUser").(model.User)

	reminds, nextCursor, err := server.TodoStorage.GetCompletedReminds(server.ctx, fetchParams, user.ID)

	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	res := model.TodoResponse{
		Todos: reminds,
		PageInfo: pagination.PageInfo{
			Page:       fetchParams.Page,
			NextCursor: nextCursor,
		},
	}

	utils.JSONFormat(w, http.StatusOK, res)
}

// GetAllReminds makes request to DB for all reminds. Works with cursor pagination
// GetAllReminds handle get completed reminds.
//
//	@Description	GetAllReminds
//	@Summary		return a list of all reminds
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		string	true	"limit"
//	@Param			cursor	query		string	true	"cursor"
//	@Success		200		{object}	model.TodoResponse
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Security		BasicAuth
//	@Router			/remind [get]
func (server *Server) GetAllReminds(w http.ResponseWriter, r *http.Request) {
	// scan for limit in parameters
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("limit parameter is invalid, should be positive integer"))
		return
	}

	// by default limit = 5
	if limit == 0 {
		limit = 5
	}

	// scan for cursor in parameters
	cursorStr := r.URL.Query().Get("cursor")
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil && cursorStr != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("cursor parameter is invalid"))
		return
	}

	//initialize fetchParameters
	fetchParams := pagination.Page{
		Limit:  limit,
		Cursor: cursor,
	}

	user := r.Context().Value("currentUser").(model.User)

	var reminds []model.Todo
	var nextCursor int

	cacheRemindsKey := fmt.Sprintf("All reminds for user:%d", user.ID)
	cacheCursorKey := "All reminds cursor"
	exist, _ := server.cache.IfExistsInCache(server.ctx, cacheRemindsKey)
	if exist {
		cachedReminds, err := server.cache.Get(server.ctx, cacheRemindsKey)
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}

		unmarshalErr := json.Unmarshal([]byte(cachedReminds), &reminds)
		if unmarshalErr != nil {
			utils.JSONError(w, http.StatusInternalServerError, unmarshalErr)
			return
		}

		cachedCursor, err := server.cache.Get(server.ctx, cacheCursorKey)
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}

		nextCursor, _ = strconv.Atoi(cachedCursor)

		server.Logger.Println("Serving All reminds from cache ...")
	} else {
		reminds, nextCursor, err = server.TodoStorage.GetAllReminds(server.ctx, fetchParams, user.ID)
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}

		byteUsers, marshalErr := json.Marshal(reminds)
		if marshalErr == nil {
			err := server.cache.Set(server.ctx, cacheRemindsKey, string(byteUsers))
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}
		}

		err := server.cache.Set(server.ctx, cacheCursorKey, strconv.Itoa(nextCursor))
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

	res := model.TodoResponse{
		Todos: reminds,
		PageInfo: pagination.PageInfo{
			Page:       fetchParams,
			NextCursor: nextCursor,
		},
	}

	utils.JSONFormat(w, http.StatusOK, res)
}

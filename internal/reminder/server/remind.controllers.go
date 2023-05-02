package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
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

// TODO: why you choose this way as a docs?

// AddRemind godoc
//
//	@Description	AddRemind
//	@Summary		create a new remind
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			input	body		domain.TodoInput	true	"remind info"
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

	if input.Description == "" || input.DeadlineAt == "" || input.Title == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("nothing to save"))
		return
	}

	userID := r.Context().Value("userID").(string)

	var todo model.Todo

	deadlineParseTime, err := time.Parse(time.RFC3339, input.DeadlineAt)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	createParseTime, err := time.Parse("02.01.2006, 15:04:05", input.CreatedAt)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	np := make([]time.Time, len(input.NotifyPeriod))

	if len(input.NotifyPeriod) > 0 {
		for _, period := range input.NotifyPeriod {
			periodParseTime, err := time.Parse(time.RFC3339, period)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, err)
				return
			}
			if periodParseTime.After(deadlineParseTime) {
				utils.JSONError(w, http.StatusBadRequest, errors.New("time to deadline notification can't be more than deadline time"))
				return
			}
			if periodParseTime.Before(deadlineParseTime.AddDate(0, 0, -2)) {
				utils.JSONError(w, http.StatusBadRequest, errors.New("time to deadline notification can't be less than 2 days to deadline time"))
				return
			}
			np = append(np, periodParseTime.Truncate(time.Minute))
		}
	}

	todo.CreatedAt = createParseTime
	todo.Description = input.Description
	todo.Title = input.Title
	todo.DeadlineAt = deadlineParseTime.Truncate(time.Minute)
	todo.UserID = userID
	todo.DeadlineNotify = input.DeadlineNotify
	todo.NotifyPeriod = np

	remind, err := server.TodoStorage.CreateRemind(server.ctx, todo)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusCreated, remind)
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
//	@Success		200		{object}	domain.TodoResponse
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

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("filter parameter is invalid"))
		return
	}

	FilterOption := r.URL.Query().Get("filterOption")
	if FilterOption == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("FilterOption parameter is invalid"))
		return
	}

	fetchParam := pagination.Page{
		Limit:        limit,
		Cursor:       cursor,
		Filter:       filter,
		FilterOption: FilterOption,
	}

	userID := r.Context().Value("userID").(string)

	reminds, totalCount, nextCursor, err := server.TodoStorage.GetNewReminds(server.ctx, fetchParam, userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	res := model.TodoResponse{
		Todos: reminds,
		Count: totalCount,
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
//	@Success		200	{object}	domain.Todo
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
//	@Param			input	body		domain.TodoUpdateInput	true	"update info"
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

	if input.Title == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("description is empty"))
		return
	}

	remind, err := server.TodoStorage.UpdateRemind(server.ctx, rID, input)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, remind)
}

func (server *Server) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uID := vars["id"]

	var input model.UserConfigs

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.TodoStorage.UpdateUserConfig(server.ctx, uID, input)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, "success")
}

// UpdateCompleteStatus update Completed field to true
//
//	@Description	UpdateCompleteStatus
//	@Summary		update remind's field "completed"
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			input	body		domain.TodoUpdateStatusInput	true	"update info"
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
//	@Success		200		{object}	domain.TodoResponse
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

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("filter parameter is invalid"))
		return
	}

	FilterOption := r.URL.Query().Get("filterOption")
	if FilterOption == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("FilterOption parameter is invalid"))
		return
	}

	//initialize fetchParameters
	fetchParams := storage.Params{
		Page: pagination.Page{
			Cursor:       cursor,
			Limit:        limit,
			FilterOption: FilterOption,
			Filter:       filter,
		},
		TimeRangeFilter: storage.TimeRangeFilter{
			StartRange: rangeStart,
			EndRange:   rangeEnd,
		},
	}

	userID := r.Context().Value("userID").(string)

	reminds, count, nextCursor, err := server.TodoStorage.GetCompletedReminds(server.ctx, fetchParams, userID)

	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	res := model.TodoResponse{
		Todos: reminds,
		Count: count,
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
//	@Success		200		{object}	domain.TodoResponse
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

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("filter parameter is invalid"))
		return
	}

	FilterOption := r.URL.Query().Get("filterOption")
	if FilterOption == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("FilterOption parameter is invalid"))
		return
	}

	//initialize fetchParameters
	fetchParams := pagination.Page{
		Limit:        limit,
		Cursor:       cursor,
		Filter:       filter,
		FilterOption: FilterOption,
	}

	userID := r.Context().Value("userID").(string)

	reminds, count, nextCursor, err := server.TodoStorage.GetAllReminds(server.ctx, fetchParams, userID)

	res := model.TodoResponse{
		Todos: reminds,
		Count: count,
		PageInfo: pagination.PageInfo{
			Page:       fetchParams,
			NextCursor: nextCursor,
		},
	}

	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, res)
}

func (server *Server) GetOrCreateUserConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uID := vars["id"]

	if uID == "" {
		utils.JSONError(w, http.StatusInternalServerError, errors.New("empty or wrong userID"))
		return
	}

	var userConfigs model.UserConfigs
	var err error

	userConfigs, err = server.TodoStorage.GetUserConfigs(server.ctx, uID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	} else if userConfigs == (model.UserConfigs{}) {
		userConfigs, err = server.TodoStorage.CreateUserConfigs(server.ctx, uID)
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.JSONFormat(w, http.StatusOK, userConfigs)
}

func (server *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		authorizationHeader := r.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("you are not logged in"))
			return
		}

		verifyToken, err := server.FireClient.VerifyIDToken(server.ctx, token)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("error verify token"))
			return
		}

		userID := verifyToken.Claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

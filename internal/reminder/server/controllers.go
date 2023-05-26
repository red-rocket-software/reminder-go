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
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/utils"
)

// AddRemind
// @Description	AddRemind
// @Summary		create a new remind
// @Tags			reminds
// @Accept			json
// @Produce		json
// @Param			input	body		domain.TodoInput	true	"remind info"
// @Success		201		{string}	domain.Todo
//
// @Failure		422		{object}	utils.HTTPError
// @Failure		400		{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
//
// @Router			/remind [post]
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

// DeleteRemind
// @Description	DeleteRemind
// @Summary		delete remind
// @Tags			reminds
// @Accept			json
// @Produce		json
// @Param			id	path		int		true	"id"
// @Success		204	{string}	string	"remind with id:1 successfully deleted"
//
// @Failure		400	{object}	utils.HTTPError
// @Failure		500	{object}	utils.HTTPError
//
// @Router			/remind{id} [delete]
func (server *Server) DeleteRemind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	remindID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	// deleting remind from db
	if err := server.TodoStorage.DeleteRemind(server.ctx, remindID); err != nil {
		if errors.Is(err, model.ErrDeleteFailed) {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}
		if errors.Is(err, model.ErrCantFindRemindWithID) {
			utils.JSONError(w, http.StatusNotFound, err)
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	successMsg := fmt.Sprintf("remind with id:%d successfully deleted", remindID)

	w.Header().Set("Content-Type", "application/json")
	utils.JSONFormat(w, http.StatusNoContent, successMsg)
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
//	@Success		200		{string}	domain.Todo
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
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

// UpdateUserConfig update user_config model
//
//	@Description	UpdateUserConfig
//	@Summary		update user_config with given fields
//	@Tags			user_config
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"id"
//	@Param			input	body		domain.UserConfigs	true	"update info"
//	@Success		200		{string}	string					"success"
//
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Router			/configs/{id} [put]
func (server *Server) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uID := vars["id"]

	var input model.UserConfigs

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.ConfigsStorage.UpdateUserConfig(server.ctx, uID, input)
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

	tn := time.Now().Truncate(1 * time.Second)

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

// GetReminds handle get reminds.
//
//	@Description	GetReminds
//	@Summary		return a list of reminds according to params
//	@Tags			reminds
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		string	true	"limit"
//	@Param			cursor	query		string	true	"cursor"
//	@Param			filter	query		string	true	"filter"
//	@Param			filterOptions	query		string	true	"filterOptions"
//	@Success		200		{object}	domain.TodoResponse
//
//	@Failure		400		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//
//	@Router			/remind [get]
func (server *Server) GetReminds(w http.ResponseWriter, r *http.Request) {
	// scan for limit in parameters
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("limit parameter is invalid, should be positive integer"))
		return
	}

	// by default limit = 10
	if limit == 0 {
		limit = 10
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

	filterOption := r.URL.Query().Get("filterOption")
	if filterOption == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("FilterOption parameter is invalid"))
		return
	}

	filterParams := r.URL.Query().Get("filterParams")
	if filterParams == "" {
		utils.JSONError(w, http.StatusBadRequest, errors.New("filterParams parameter is invalid"))
		return
	}

	//initialize fetchParameters
	params := model.FetchParams{
		Page: utils.Page{
			Cursor: cursor,
			Limit:  limit,
		},
		FilterByDate:  filter,
		FilterBySort:  filterOption,
		FilterByQuery: filterParams,
	}

	userID := r.Context().Value("userID").(string)

	reminds, count, nextCursor, err := server.TodoStorage.GetReminds(server.ctx, params, userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	res := model.TodoResponse{
		Todos: reminds,
		Count: count,
		PageInfo: utils.PageInfo{
			Page:       params.Page,
			NextCursor: nextCursor,
		},
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

	userConfigs, err = server.ConfigsStorage.GetUserConfigs(server.ctx, uID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	} else if userConfigs == (model.UserConfigs{}) {
		userConfigs, err = server.ConfigsStorage.CreateUserConfigs(server.ctx, uID)
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.JSONFormat(w, http.StatusOK, userConfigs)
}

func (server *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONFormat(w, http.StatusOK, "OK")
}

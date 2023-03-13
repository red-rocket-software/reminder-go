package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
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

	cacheKey := r.URL.String()
	byteUser, marshalErr := json.Marshal(currentUser)
	if marshalErr == nil {
		err := server.cache.Set(server.ctx, cacheKey, string(byteUser))
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

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

	err = server.TodoStorage.UpdateUserNotification(server.ctx, uID, input.Notification)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusOK, "user notification status successfully updated")
}

func (s *Server) cacheFetchMeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey := r.URL.String()

		var user model.User

		exist, err := s.cache.IfExistsInCache(s.ctx, cacheKey)
		if err != nil {
			return
		}
		if exist {
			result, err := s.cache.Get(s.ctx, cacheKey)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}
			unmarshalErr := json.Unmarshal([]byte(result), &user)
			if unmarshalErr != nil {
				utils.JSONError(w, http.StatusInternalServerError, unmarshalErr)
				return
			}
			s.Logger.Println("Serving fetchMe from cache ...")
			utils.JSONFormat(w, http.StatusOK, user)
			return
		}

		next.ServeHTTP(w, r)
	}
}

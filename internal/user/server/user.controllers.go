package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
	model "github.com/red-rocket-software/reminder-go/internal/user/domain"
	"github.com/red-rocket-software/reminder-go/utils"
)

func (server *Server) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	var pathURL = "/"

	if r.URL.Query().Get("state") != "" {
		pathURL = r.URL.Query().Get("state")
	}

	if code == "" {
		utils.JSONError(w, http.StatusUnauthorized, errors.New("authorization code not provided"))
		return
	}

	tokenRes, err := utils.GetGoogleOuathToken(code, server.config)
	if err != nil {
		utils.JSONError(w, http.StatusBadGateway, err)
		return
	}

	googleUser, err := utils.GetGoogleUser(tokenRes.AccessToken, tokenRes.IDToken)
	if err != nil {
		utils.JSONError(w, http.StatusBadGateway, err)
		return
	}

	now := time.Now()
	email := strings.ToLower(googleUser.Email)

	dataUser := model.User{
		Name:      googleUser.Name,
		Email:     email,
		Password:  "",
		Photo:     googleUser.Picture,
		Provider:  "Google",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := server.UserStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			_, err = server.UserStorage.CreateUser(server.ctx, dataUser)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, err)
				return
			}
		}
	}

	user, err = server.UserStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	accessToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtSecret, server.config.Auth.TokenExpiredIn)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	refreshToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtRefreshSecret, server.config.Auth.JwtRefreshKeyExpire)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	accessTokenCookie := utils.CreateCookie("token", accessToken, "/", "localhost", server.config.Auth.TokenMaxAge*60, false, true)
	refreshTokenCookie := utils.CreateCookie("refresh_token", refreshToken, "/", "localhost", server.config.Auth.RefreshTokenMaxAge*60, false, true)

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	http.Redirect(w, r, fmt.Sprint(server.config.Auth.FrontendOrigin, pathURL), http.StatusTemporaryRedirect)
}

// SignUpUser godoc
//
//	@Summary		SignUpUser
//	@Tags			user
//	@Description	create user account
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		domain.RegisterUserInput	true	"user info"
//	@Success		201		{string}	string					"User is successfully created id: 1"
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//	@Router			/register [post]
func (server *Server) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload model.RegisterUserInput

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if len(payload.Password) < 6 {
		utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("password should be more then 6"))
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	now := time.Now()

	newUser := model.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Provider:  "local",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = model.Validate(newUser, "")
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	id, err := server.UserStorage.CreateUser(server.ctx, newUser)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("user with this email is already existing"))
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONFormat(w, http.StatusCreated, fmt.Sprintf("User is successfully created id:%d", id))
}

// SignInUser godoc
//
//	@Summary		SignInUser
//	@Tags			user
//	@Description	user user, return user and save token to cookie
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		domain.LoginUserInput	true	"user email and password"
//	@Success		201		{object}	domain.User
//	@Failure		422		{object}	utils.HTTPError
//	@Failure		401		{object}	utils.HTTPError
//	@Failure		500		{object}	utils.HTTPError
//	@Router			/login [post]
func (server *Server) SignInUser(w http.ResponseWriter, r *http.Request) {
	var payload model.LoginUserInput

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := server.UserStorage.GetUserByEmail(server.ctx, payload.Email)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, fmt.Errorf("invalid email or Password %v", err))
		return
	}

	err = utils.CheckPassword(payload.Password, user.Password)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, err)
		return
	}

	if user.Provider == "Google" {
		utils.JSONFormat(w, http.StatusUnauthorized, fmt.Sprintf("Use %v OAuth instead", user.Provider))
		return
	}

	accessToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtSecret, server.config.Auth.TokenExpiredIn)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	refreshToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtRefreshSecret, server.config.Auth.JwtRefreshKeyExpire)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	accessTokenCookie := utils.CreateCookie("token", accessToken, "/", "localhost", server.config.Auth.TokenMaxAge*60, false, true)
	refreshTokenCookie := utils.CreateCookie("refresh_token", refreshToken, "/", "localhost", server.config.Auth.RefreshTokenMaxAge*60, false, true)

	userResponse := model.ToResponseUser(user)

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	utils.JSONFormat(w, http.StatusCreated, userResponse)
}

// LogOutUser godoc
//
//	@Summary		LogOutUser
//	@Tags			user
//	@Description	logout user and remove cookie
//	@Produce		json
//	@Success		200	{string}	string "Success"
//
//	@Security		BasicAuth
//
//	@Router			/logout [get]
func (server *Server) LogOutUser(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie := utils.CreateExpiredCookie("token", "", "/", "localhost", false, true)
	refreshTokenCookie := utils.CreateExpiredCookie("refresh_token", "", "/", "localhost", false, true)

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	utils.JSONFormat(w, http.StatusOK, "Success")
}

func (server *Server) SignInOrSignUp(w http.ResponseWriter, r *http.Request) {
	var payload model.LoginUserInput

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := server.UserStorage.GetUserByEmail(server.ctx, payload.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {

			hashedPassword, err := utils.HashPassword(payload.Password)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}

			now := time.Now()

			newUser := model.User{
				Name:      "",
				Email:     strings.ToLower(payload.Email),
				Password:  hashedPassword,
				Provider:  "local",
				Verified:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = model.Validate(newUser, "login")
			if err != nil {
				utils.JSONError(w, http.StatusUnprocessableEntity, err)
				return
			}

			_, err = server.UserStorage.CreateUser(server.ctx, newUser)
			if err != nil {
				if strings.Contains(err.Error(), "23505") {
					utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("user with this email is already existing"))
					return
				}
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}

			createdUser, err := server.UserStorage.GetUserByEmail(server.ctx, payload.Email)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, fmt.Errorf("invalid email or Password %v", err))
				return
			}

			accessToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtSecret, server.config.Auth.TokenExpiredIn)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}
			refreshToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtRefreshSecret, server.config.Auth.JwtRefreshKeyExpire)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}

			accessTokenCookie := utils.CreateCookie("token", accessToken, "/", "localhost", server.config.Auth.TokenMaxAge*60, false, true)
			refreshTokenCookie := utils.CreateCookie("refresh_token", refreshToken, "/", "localhost", server.config.Auth.RefreshTokenMaxAge*60, false, true)

			http.SetCookie(w, &accessTokenCookie)
			http.SetCookie(w, &refreshTokenCookie)

			resUser := model.ToResponseUser(createdUser)

			utils.JSONFormat(w, http.StatusCreated, resUser)
			return
		} else {
			utils.JSONError(w, http.StatusBadRequest, fmt.Errorf("invalid email or Password %v", err))
			return
		}
	}

	accessToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtSecret, server.config.Auth.TokenExpiredIn)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	refreshToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtRefreshSecret, server.config.Auth.JwtRefreshKeyExpire)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	accessTokenCookie := utils.CreateCookie("token", accessToken, "/", "localhost", server.config.Auth.TokenMaxAge*60, false, true)
	refreshTokenCookie := utils.CreateCookie("refresh_token", refreshToken, "/", "localhost", server.config.Auth.RefreshTokenMaxAge*60, false, true)

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	resUser := model.ToResponseUser(user)

	utils.JSONFormat(w, http.StatusCreated, resUser)
}

func (server *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie("token")

		authorizationHeader := r.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie.Value
		}

		if token == "" {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("you are not logged in"))
			return
		}

		sub, err := utils.ValidateToken(token, server.config.Auth.JwtSecret)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, err)
			return
		}

		user, err := server.UserStorage.GetUserByID(server.ctx, int(sub.(float64)))
		if err != nil {
			utils.JSONError(w, http.StatusBadRequest, err)
		}

		ctx := context.WithValue(r.Context(), "currentUser", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetMe godoc
//
//	@Description	GetMe
//	@Summary		fetch current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	domain.User
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
//	@Param			input	body		domain.NotificationUserInput	true	"update info"
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

	err = server.UserStorage.UpdateUserNotification(server.ctx, uID, input)
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
	if err := server.UserStorage.DeleteUser(server.ctx, userID); err != nil {
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

	accessTokenCookie := utils.CreateExpiredCookie("token", "", "/", "localhost", false, true)
	refreshTokenCookie := utils.CreateExpiredCookie("refresh_token", "", "/", "localhost", false, true)

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	successMsg := fmt.Sprintf("user with id:%d successfully deleted", userID)

	w.Header().Set("Content-Type", "application/json")
	utils.JSONFormat(w, http.StatusNoContent, successMsg)
}

func (server *Server) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	message := "could not refresh access token"

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, errors.New(message))
		return
	}

	sub, err := utils.ValidateToken(cookie.Value, server.config.Auth.JwtRefreshSecret)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err)
		return
	}

	user, err := server.UserStorage.GetUserByID(server.ctx, int(sub.(float64)))
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, errors.New("the user belonging to this token no logger exists"))
		return
	}

	accessToken, err := utils.GenerateNewToken(user.ID, server.config.Auth.JwtSecret, server.config.Auth.TokenExpiredIn)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err)
		return
	}

	accessTokenCookie := utils.CreateCookie("token", accessToken, "/", "localhost", server.config.Auth.TokenMaxAge*60, false, true)

	http.SetCookie(w, &accessTokenCookie)

	utils.JSONFormat(w, http.StatusOK, accessToken)
}

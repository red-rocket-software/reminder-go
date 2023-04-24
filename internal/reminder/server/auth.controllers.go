package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	model2 "github.com/red-rocket-software/reminder-go/internal/reminder/app/domain"
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

	dataUser := model2.User{
		Name:      googleUser.Name,
		Email:     email,
		Password:  "",
		Photo:     googleUser.Picture,
		Provider:  "Google",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			_, err = server.TodoStorage.CreateUser(server.ctx, dataUser)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, err)
				return
			}
		}
	}

	user, err = server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, user.ID, server.config.Auth.JwtSecret)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = server.config.Auth.TokenMaxAge * 60
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, fmt.Sprint(server.config.Auth.FrontendOrigin, pathURL), http.StatusTemporaryRedirect)
}

func (server *Server) GithubAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	var pathURL = "/"

	if r.URL.Query().Get("state") != "" {
		pathURL = r.URL.Query().Get("state")
	}

	if code == "" {
		utils.JSONError(w, http.StatusUnauthorized, errors.New("authorization code not provided"))
		return
	}

	tokenRes, err := utils.GetGithubOuathToken(code, server.config)
	if err != nil {
		utils.JSONError(w, http.StatusBadGateway, err)
		return
	}

	githubUser, err := utils.GetGithubUser(tokenRes)

	if err != nil {
		utils.JSONError(w, http.StatusBadGateway, err)
		return
	}

	now := time.Now()
	email := strings.ToLower(githubUser.Email)

	dataUser := model2.User{
		Name:      githubUser.Name,
		Email:     email,
		Password:  "",
		Provider:  "Github",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			_, err = server.TodoStorage.CreateUser(server.ctx, dataUser)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, err)
				return
			}
		}
	}

	user, err = server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, user.ID, server.config.Auth.JwtSecret)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = server.config.Auth.TokenMaxAge * 60
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, fmt.Sprint(server.config.Auth.FrontendOrigin, pathURL), http.StatusTemporaryRedirect)
}

func (server *Server) LinkedinAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	var pathURL = "/"

	if r.URL.Query().Get("state") != "" {
		pathURL = r.URL.Query().Get("state")
	}

	if code == "" {
		utils.JSONError(w, http.StatusUnauthorized, errors.New("authorization code not provided"))
		return
	}

	tokenRes, err := utils.GetLinkedinOauthToken(code, server.config)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	linkedinUser, err := utils.GetLinkedinUser(server.config, tokenRes)

	if err != nil {
		utils.JSONError(w, http.StatusBadGateway, err)
		return
	}

	now := time.Now()
	email := strings.ToLower(linkedinUser.Email)

	dataUser := model2.User{
		Name:      linkedinUser.FirstName,
		Email:     email,
		Password:  "",
		Photo:     linkedinUser.Picture,
		Provider:  "Linkedin",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			_, err = server.TodoStorage.CreateUser(server.ctx, dataUser)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, err)
				return
			}
		}
	}

	user, err = server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, user.ID, server.config.Auth.JwtSecret)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = server.config.Auth.TokenMaxAge * 60
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

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

	newUser := model2.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Provider:  "local",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = model2.Validate(newUser, "")
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	id, err := server.TodoStorage.CreateUser(server.ctx, newUser)
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

	user, err := server.TodoStorage.GetUserByEmail(server.ctx, payload.Email)
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

	token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, user.ID, server.config.Auth.JwtSecret)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.ToResponseUser(user)

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = server.config.Auth.TokenMaxAge
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

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
	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.MaxAge = -1
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

	utils.JSONFormat(w, http.StatusOK, "Success")
}

func (server *Server) SignInOrSignUp(w http.ResponseWriter, r *http.Request) {
	var payload model.LoginUserInput

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := server.TodoStorage.GetUserByEmail(server.ctx, payload.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {

			hashedPassword, err := utils.HashPassword(payload.Password)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}

			now := time.Now()

			newUser := model2.User{
				Name:      "",
				Email:     strings.ToLower(payload.Email),
				Password:  hashedPassword,
				Provider:  "local",
				Verified:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = model2.Validate(newUser, "login")
			if err != nil {
				utils.JSONError(w, http.StatusUnprocessableEntity, err)
				return
			}

			_, err = server.TodoStorage.CreateUser(server.ctx, newUser)
			if err != nil {
				if strings.Contains(err.Error(), "23505") {
					utils.JSONError(w, http.StatusUnprocessableEntity, errors.New("user with this email is already existing"))
					return
				}
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}

			createdUser, err := server.TodoStorage.GetUserByEmail(server.ctx, payload.Email)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, fmt.Errorf("invalid email or Password %v", err))
				return
			}

			token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, createdUser.ID, server.config.Auth.JwtSecret)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err)
				return
			}

			cookie := http.Cookie{}
			cookie.Name = "token"
			cookie.Value = token
			cookie.Path = "/"
			cookie.Domain = "localhost"
			cookie.Secure = false
			cookie.HttpOnly = true
			http.SetCookie(w, &cookie)

			resUser := model.ToResponseUser(createdUser)

			utils.JSONFormat(w, http.StatusCreated, resUser)
			return
		} else {
			utils.JSONError(w, http.StatusBadRequest, fmt.Errorf("invalid email or Password %v", err))
			return
		}
	}

	token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, user.ID, server.config.Auth.JwtSecret)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.Secure = false
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)

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

		user, err := server.TodoStorage.GetUserByID(server.ctx, int(sub.(float64)))
		if err != nil {
			utils.JSONError(w, http.StatusBadRequest, err)
		}

		ctx := context.WithValue(r.Context(), "currentUser", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

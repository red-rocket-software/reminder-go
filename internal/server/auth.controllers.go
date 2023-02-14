package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/utils"
)

func (server *Server) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	var pathURL string = "/"

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

	googleUser, err := utils.GetGoogleUser(tokenRes.AccessToken, tokenRes.IdToken)
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
		Provider:  "Google",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email); err != nil {
		server.TodoStorage.CreateUser(server.ctx, dataUser)
	}

	user := server.TodoStorage.GetUserByEmail(server.ctx, dataUser.Email)

	token, err := utils.GenerateToken(server.config.Auth.TokenExpiredIn, user.Id.String(), server.config.Auth.JwtSecret)
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

	http.Redirect(w, r, server.config.Auth.FrontendOrigin, http.StatusTemporaryRedirect)
}

func (server *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("error occurred while reading cookie"))
		}

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

		user := server.TodoStorage.GetUserByEmail(server.ctx, fmt.Sprint(sub))
		ctx := context.WithValue(r.Context(), "currentUser", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

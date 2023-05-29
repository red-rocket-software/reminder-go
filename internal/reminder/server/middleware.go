package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/red-rocket-software/reminder-go/pkg/utils"
)

type middlewareData struct {
	uID  string
	role string
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

		role, fireToken, err := utils.ParseToken(token, server.config.JWTSecret)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("error to parse JWT token"))
			return
		}

		verifyToken, err := server.FireClient.VerifyIDToken(fireToken)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("error verify token"))
			return
		}

		userID := verifyToken.Claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), "authData", middlewareData{
			uID:  userID,
			role: role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) RoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareData := r.Context().Value("authData").(middlewareData)
		route := strings.Split(r.URL.Path, "/")[1]

		routes, err := server.TodoStorage.GetUserRoutes(r.Context(), middlewareData.role)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, err)
			return
		}

		for _, rt := range routes {
			if rt == "all" || rt == route {
				ctx := context.WithValue(r.Context(), "userID", middlewareData.uID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		utils.JSONError(w, http.StatusUnauthorized, errors.New("you have no access"))
	})
}

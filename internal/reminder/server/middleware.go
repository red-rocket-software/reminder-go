package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
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

		role, uid, err := utils.ParseToken(token, server.config.JWTSecret)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), "authData", middlewareData{
			uID:  uid,
			role: role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) RoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareData := r.Context().Value("authData").(middlewareData)

		p, err := storage.GetUserPermissions(r.Context(), middlewareData.role, server.config)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, err)
			return
		}
		fmt.Println("permissions:", p)

		for _, rt := range p {
			fmt.Println("sub-features:", rt)

			if rt == "allReminder" {
				ctx := context.WithValue(r.Context(), "userID", middlewareData.uID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		utils.JSONError(w, http.StatusUnauthorized, errors.New("you have no access"))
	})
}

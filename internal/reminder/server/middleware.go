package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/red-rocket-software/reminder-go/pkg/utils"
)

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

		verifyToken, err := server.FireClient.VerifyIDToken(token)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, errors.New("error verify token"))
			return
		}

		userID := verifyToken.Claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

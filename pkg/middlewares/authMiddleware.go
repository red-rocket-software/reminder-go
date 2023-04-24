package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/user/storage"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
	"github.com/red-rocket-software/reminder-go/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	config := config.GetConfig()
	logger := logging.GetLogger()
	postgresClient, err := postgresql.NewClient(context.Background(), 5, *config)
	if err != nil {
		logger.Fatalf("Error create new db client:%v\n", err)
	}
	defer postgresClient.Close()
	storage := storage.NewUserStorage(postgresClient, &logger)
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

		sub, err := utils.ValidateToken(token, config.Auth.JwtSecret)
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, err)
			return
		}

		user, err := storage.GetUserByID(context.Background(), int(sub.(float64)))
		if err != nil {
			utils.JSONError(w, http.StatusBadRequest, err)
		}

		ctx := context.WithValue(r.Context(), "currentUser", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

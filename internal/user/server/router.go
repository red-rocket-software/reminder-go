package server

import (
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/pkg/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

// ConfigureAuthRouter returns router with routes from controllers
func (server *Server) ConfigureAuthRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.CorsMiddleware)

	// private routes
	privateRoute := router.PathPrefix("").Subrouter()
	privateRoute.Use(server.AuthMiddleware)

	privateRoute.HandleFunc("/logout", server.LogOutUser).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/fetchMe", server.GetMe).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/user/{id}", server.DeleteUser).Methods("DELETE", "OPTIONS")

	// public routes
	publicRoute := router.PathPrefix("").Subrouter()

	// login handlers
	publicRoute.HandleFunc("/register", server.SignUpUser).Methods("POST", "OPTIONS")
	publicRoute.HandleFunc("/login", server.SignInUser).Methods("POST", "OPTIONS")
	publicRoute.HandleFunc("/login-or-register", server.SignInOrSignUp).Methods("POST", "OPTIONS")

	// login callbacks
	publicRoute.HandleFunc("/google/callback", server.GoogleAuth).Methods("GET", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}

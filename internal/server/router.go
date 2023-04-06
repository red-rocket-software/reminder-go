package server

import (
	"github.com/gorilla/mux"
	_ "github.com/red-rocket-software/reminder-go/docs"
	"github.com/red-rocket-software/reminder-go/pkg/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

// ConfigureRouter returns router with routes from controllers
func (server *Server) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.CorsMiddleware)

	// private routes
	privateRoute := router.PathPrefix("").Subrouter()
	privateRoute.Use(server.AuthMiddleware)

	privateRoute.HandleFunc("/remind", server.GetAllReminds).Methods("GET")
	privateRoute.HandleFunc("/remind/{id}", server.GetRemindByID).Methods("GET")
	privateRoute.HandleFunc("/remind", server.AddRemind).Methods("POST", "OPTIONS")
	privateRoute.HandleFunc("/remind/{id}", server.UpdateRemind).Methods("PUT")
	privateRoute.HandleFunc("/status/{id}", server.UpdateCompleteStatus).Methods("PUT", "OPTIONS")
	privateRoute.HandleFunc("/remind/{id}", server.DeleteRemind).Methods("DELETE", "OPTIONS")
	privateRoute.HandleFunc("/completed", server.GetCompletedReminds).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/current", server.GetCurrentReminds).Methods("GET", "OPTIONS")

	privateRoute.HandleFunc("/logout", server.LogOutUser).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/fetchMe", server.GetMe).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/user/{id}", server.UpdateUserNotification).Methods("PUT", "OPTIONS")
	privateRoute.HandleFunc("/user/{id}", server.DeleteUser).Methods("DELETE", "OPTIONS")

	// public routes
	publicRoute := router.PathPrefix("").Subrouter()

	// login handlers
	publicRoute.HandleFunc("/register", server.SignUpUser).Methods("POST", "OPTIONS")
	publicRoute.HandleFunc("/login", server.SignInUser).Methods("POST", "OPTIONS")
	publicRoute.HandleFunc("/login-or-register", server.SignInOrSignUp).Methods("POST", "OPTIONS")

	// login callbacks
	publicRoute.HandleFunc("/google/callback", server.GoogleAuth).Methods("GET", "OPTIONS")
	publicRoute.HandleFunc("/linkedin/callback", server.LinkedinAuth).Methods("GET", "OPTIONS")
	publicRoute.HandleFunc("/github/callback", server.GithubAuth).Methods("GET", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}

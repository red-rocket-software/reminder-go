package server

import (
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/pkg/middlewares"
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
	privateRoute.HandleFunc("/completed", server.GetCompletedReminds).Methods("GET")
	privateRoute.HandleFunc("/current", server.GetCurrentReminds).Methods("GET")

	privateRoute.HandleFunc("/logout", server.LogOutUser).Methods("GET", "OPTIONS")

	// public routes
	publicRoute := router.PathPrefix("").Subrouter()

	publicRoute.HandleFunc("/register", server.SignUpUser).Methods("POST", "OPTIONS")
	publicRoute.HandleFunc("/login", server.SignInUser).Methods("POST", "OPTIONS")

	publicRoute.HandleFunc("/google/callback", server.GoogleAuth).Methods("GET")

	return router
}

package server

import (
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/pkg/middlewares"
)

// ConfigureRouter returns router with routes from controllers
func (server *Server) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.CorsMiddleware)

	router.HandleFunc("/remind", server.GetAllReminds).Methods("GET")
	router.HandleFunc("/remind/{id}", server.GetRemindByID).Methods("GET")
	router.HandleFunc("/remind", server.AuthMiddleware(server.AddRemind)).Methods("POST", "OPTIONS")
	router.HandleFunc("/remind/{id}", server.UpdateRemind).Methods("PUT")
	router.HandleFunc("/status/{id}", server.UpdateCompleteStatus).Methods("PUT", "OPTIONS")
	router.HandleFunc("/remind/{id}", server.DeleteRemind).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/completed", server.GetCompletedReminds).Methods("GET")
	router.HandleFunc("/current", server.GetCurrentReminds).Methods("GET")

	router.HandleFunc("/auth/google/callback", server.GoogleAuth).Methods("GET")

	authGroup := router.PathPrefix("/auth").Subrouter()
	authGroup.HandleFunc("/register", server.SignUpUser).Methods("POST", "OPTIONS")
	authGroup.HandleFunc("/login", server.SignInUser).Methods("POST", "OPTIONS")
	authGroup.HandleFunc("/logout", server.AuthMiddleware(server.LogOutUser)).Methods("GET", "OPTIONS")

	return router

}

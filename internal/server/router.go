package server

import (
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/pkg/middlewares"
)

// function ConfigureRouter returns router with routes from controllers
func (server *Server) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.CorsMiddleware)

	router.HandleFunc("/remind", server.GetAllReminds).Methods("GET")
	router.HandleFunc("/remind/{id}", server.GetRemindByID).Methods("GET")
	router.HandleFunc("/remind", server.AddRemind).Methods("POST", "OPTIONS")
	router.HandleFunc("/remind/{id}", server.UpdateRemind).Methods("PUT")
	router.HandleFunc("/status/{id}", server.UpdateCompleteStatus).Methods("PUT", "OPTIONS")
	router.HandleFunc("/remind/{id}", server.DeleteRemind).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/completed", server.GetCompletedReminds).Methods("GET")
	router.HandleFunc("/current", server.GetCurrentReminds).Methods("GET")

	return router

}

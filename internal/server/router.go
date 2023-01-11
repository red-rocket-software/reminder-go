package server

import (
	"github.com/gorilla/mux"
)

// function ConfigureRouter returns router with routes from controllers
func (server *Server) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	// please implement controllers methods and uncomment this rows

	router.HandleFunc("/remind", server.GetAllReminds).Methods("GET")
	router.HandleFunc("/remind/{id}", server.GetRemindByID).Methods("GET")
	router.HandleFunc("/remind", server.AddRemind).Methods("POST")
	router.HandleFunc("/remind/{id}", server.UpdateRemind).Methods("PUT")
	router.HandleFunc("/remind/{id}", server.DeleteRemind).Methods("DELETE")
	router.HandleFunc("/completed", server.GetComplitedReminds).Methods("GET")
	router.HandleFunc("/current", server.GetCurrentReminds).Methods("GET")

	return router

}

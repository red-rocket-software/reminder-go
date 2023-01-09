package server

import (
	"github.com/gorilla/mux"
)

// function ConfigureRouter returns router with routes from controllers
func (s *Server) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	// please implement controllers methods and uncomment this rows

	// router.HandleFunc("/remind", controllers.GetAllReminds).Methods("GET")
	router.HandleFunc("/remind/{id}", s.GetRemindById).Methods("GET")
	router.HandleFunc("/remind", s.AddRemind).Methods("POST")
	// router.HandleFunc("/remind/{id}", controllers.DeleteRemind).Methods("DELETE")
	router.HandleFunc("/remind/{id}", s.UpdateRemind).Methods("PUT")
	// router.HandleFunc("/completed", controllers.GetComplitedReminds).Methods("GET")
	router.HandleFunc("/current", s.GetCurrentReminds).Methods("GET")

	return router

}

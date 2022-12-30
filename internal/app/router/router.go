package router

import (
	"github.com/gorilla/mux"
	// "github.com/red-rocket-software/reminder-go/internal/app/controllers"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
)

// function ConfigureRouter returns router with routes from controllers
func ConfigureRouter(logger logging.Logger) *mux.Router {
	router := mux.NewRouter()

	// please implement controllers methods and uncomment this rows

	// router.HandleFunc("/remind", controllers.GetAllReminds(logger)).Methods("GET")
	// router.HandleFunc("/remind/{id}", controllers.GetRemindById(logger)).Methods("GET")
	// router.HandleFunc("/remind", controllers.AddRemind(logger)).Methods("POST")
	// router.HandleFunc("/remind/{id}", controllers.DeleteRemind(logger)).Methods("DELETE")
	// router.HandleFunc("/remind/{id}", controllers.UpdateRemind(logger)).Methods("PUT")
	// router.HandleFunc("/completed", controllers.GetComplitedReminds(logger)).Methods("GET")
	// router.HandleFunc("/current", controllers.GetCurrentReminds(logger)).Methods("GET")

	return router

}

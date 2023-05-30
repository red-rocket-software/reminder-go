package server

import (
	"github.com/gorilla/mux"
	_ "github.com/red-rocket-software/reminder-go/docs"
	"github.com/red-rocket-software/reminder-go/pkg/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

// ConfigureReminderRouter returns router with routes from controllers
func (server *Server) ConfigureReminderRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.Cors)

	router.HandleFunc("/health", server.HealthCheck).Methods("GET")

	// private routes
	privateRoute := router.PathPrefix("").Subrouter()
	privateRoute.Use(server.AuthMiddleware)
	privateRoute.Use(server.RoleMiddleware)

	privateRoute.HandleFunc("/reminds", server.GetReminds).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/remind/{id}", server.GetRemindByID).Methods("GET")
	privateRoute.HandleFunc("/remind", server.AddRemind).Methods("POST", "OPTIONS")
	privateRoute.HandleFunc("/remind/{id}", server.DeleteRemind).Methods("DELETE", "OPTIONS")
	privateRoute.HandleFunc("/remind/{id}", server.UpdateRemind).Methods("PUT")

	privateRoute.HandleFunc("/status/{id}", server.UpdateCompleteStatus).Methods("PUT", "OPTIONS")
	privateRoute.HandleFunc("/configs/{id}", server.GetOrCreateUserConfig).Methods("GET", "OPTIONS")
	privateRoute.HandleFunc("/configs/{id}", server.UpdateUserConfig).Methods("PUT", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}

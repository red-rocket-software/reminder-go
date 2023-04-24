package server

import (
	"github.com/gorilla/mux"
)

// ConfigureRouter returns router with routes from controllers
func (server *Server) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/test", server.test)

	//router.Use(middlewares.CorsMiddleware)
	//
	//// private routes
	//privateRoute := router.PathPrefix("").Subrouter()
	//privateRoute.Use(middlewares.AuthMiddleware)
	//
	//privateRoute.HandleFunc("/logout", server.LogOutUser).Methods("GET", "OPTIONS")
	//privateRoute.HandleFunc("/fetchMe", server.GetMe).Methods("GET", "OPTIONS")
	//privateRoute.HandleFunc("/user/{id}", server.DeleteUser).Methods("DELETE", "OPTIONS")
	//
	//// public routes
	//publicRoute := router.PathPrefix("").Subrouter()
	//
	//// login handlers
	//publicRoute.HandleFunc("/register", server.SignUpUser).Methods("POST", "OPTIONS")
	//publicRoute.HandleFunc("/login", server.SignInUser).Methods("POST", "OPTIONS")
	//publicRoute.HandleFunc("/login-or-register", server.SignInOrSignUp).Methods("POST", "OPTIONS")
	//
	//// login callbacks
	//publicRoute.HandleFunc("/google/callback", server.GoogleAuth).Methods("GET", "OPTIONS")
	//publicRoute.HandleFunc("/linkedin/callback", server.LinkedinAuth).Methods("GET", "OPTIONS")
	//publicRoute.HandleFunc("/github/callback", server.GithubAuth).Methods("GET", "OPTIONS")
	//
	//router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}

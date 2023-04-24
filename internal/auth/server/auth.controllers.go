package server

import (
	"net/http"

	"github.com/red-rocket-software/reminder-go/utils"
)

func (server *Server) test(w http.ResponseWriter, r *http.Request) {
	utils.JSONFormat(w, http.StatusOK, "HELLLLOOOOOO WORLD")
}

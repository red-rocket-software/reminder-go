package server

import (
	"net/http"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/utils"
)

func (server *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value("currentUser").(model.User)

	utils.JSONFormat(w, http.StatusOK, currentUser)
}

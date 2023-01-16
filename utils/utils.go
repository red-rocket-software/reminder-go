package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSONFormat(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func JSONError(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSONFormat(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSONFormat(w, http.StatusBadRequest, nil)
}

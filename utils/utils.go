package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonFormat(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func JsonError(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JsonFormat(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JsonFormat(w, http.StatusBadRequest, nil)
}

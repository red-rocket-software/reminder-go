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
		JSONFormat(w, statusCode, HTTPError{
			Code:    statusCode,
			Message: err.Error(),
		})
		return
	}
	JSONFormat(w, http.StatusBadRequest, nil)
}

// HTTPError response
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

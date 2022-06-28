package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GenericError struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

func WriteError(w http.ResponseWriter, code int, message string, data interface{}) {
	response := GenericError{
		Code:  code,
		Error: message,
		Data:  data,
	}

	WriteJson(w, code, response)
}

func WriteJson(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Warn("Error writing response")
	}
}

package api

import (
	v1 "finance-app/internal/api/v1"
	"finance-app/internal/database"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(db database.Database) (http.Handler, error) {
	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)
	router.PathPrefix("/api/v1").Subrouter()

	return router, nil
}

package api

import (
	"financial-app/internal/api/auth"
	v1 "financial-app/internal/api/v1"
	"financial-app/internal/database"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(db database.Database) (http.Handler, error) {
	permissions := auth.NewPermissions(db)
	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	v1.SetUserAPI(db, apiRouter, permissions)
	v1.SetUserRoleAPI(db, apiRouter, permissions)
	router.Use(auth.AutherizationToken)

	return router, nil
}

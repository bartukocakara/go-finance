package v1

import (
	"encoding/json"
	"financial-app/internal/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserAPI struct {
	DB database.Database // Will represent all database interface
}

func SetUserAPI(db database.Database, router *mux.Router) {
	api := &UserAPI{
		DB: db,
	}
	apis := []API{
		// -----------USER----------------------------
		NewAPI(http.MethodPost, "/users", api.Create),
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, func(w http.ResponseWriter, r *http.Request) {
			// an example API handler
			json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		})
	}
}

func (api *UserAPI) Create(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "user.go -> Create()")

}

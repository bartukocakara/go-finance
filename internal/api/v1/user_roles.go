package v1

import (
	"encoding/json"
	"financial-app/internal/api/auth"
	"financial-app/internal/api/utils"
	"financial-app/internal/database"
	"financial-app/internal/model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserRoleAPI struct {
	DB database.Database
}

func SetUserRoleAPI(db database.Database, router *mux.Router) {
	api := &UserRoleAPI{
		DB: db,
	}
	apis := []API{
		NewAPI(http.MethodPost, "/users/{userID}/roles", api.GrantRole),
	}
	for _, api := range apis {
		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}
}

func (api *UserRoleAPI) GrantRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user_roles.go -> GrantRole()")
	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)
	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	var userRole model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&userRole); err != nil {
		logger.WithError(err).Warn("Could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()
	if err := api.DB.GrantRole(ctx, userID, userRole.Role); err != nil {
		logger.WithError(err).Warn("Error while granting role")
		utils.WriteError(w, http.StatusInternalServerError, "Error granting role", nil)
		return
	}
	logger.WithField("UserID", userID).Info("Role Granted to user", userID)
	utils.WriteJson(w, http.StatusCreated, &ActionCreated{
		Created: true,
	})
}

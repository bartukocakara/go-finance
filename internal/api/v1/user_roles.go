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
		NewAPI(http.MethodPost, "/users/{UserID}/roles", api.GrantRole),
		NewAPI(http.MethodPut, "/users/{UserID}/roles", api.UpdateRole),
		NewAPI(http.MethodDelete, "/users/{UserID}/roles", api.RevokeRole),
		NewAPI(http.MethodGet, "/users/{UserID}/roles", api.GetRoleList),
	}
	for _, api := range apis {
		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}
}

func (api *UserRoleAPI) GrantRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user_roles.go -> GrantRole()")
	vars := mux.Vars(r)
	UserID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)
	logger = logger.WithFields(logrus.Fields{
		"userID":    UserID,
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
	if err := api.DB.GrantRole(ctx, UserID, userRole.Role); err != nil {
		logger.WithError(err).Warn("Error while granting role")
		utils.WriteError(w, http.StatusInternalServerError, "Error granting role", nil)
		return
	}
	logger.WithField("UserID", UserID).Info("Role Granted to user", UserID)
	utils.WriteJson(w, http.StatusCreated, &ActionCreated{
		Created: true,
	})
}

func (api *UserRoleAPI) UpdateRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user_role.go -> UpdateRole")
	vars := mux.Vars(r)
	UserID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger = logrus.WithFields(logrus.Fields{
		"userID":    UserID,
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
	if err := api.DB.UpdateRole(ctx, UserID, userRole.Role); err != nil {
		logger.WithError(err).Warn("Error updating role")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating role", nil)
		return
	}

	logrus.WithField("UserID", UserID).Info("User role updated successfully : ", userRole.Role)
	utils.WriteJson(w, http.StatusCreated, &ActionUpdated{
		Updated: true,
	})
}

func (api *UserRoleAPI) RevokeRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go -> Revoke Role")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
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
	if err := api.DB.RevokeRole(ctx, userID, userRole.Role); err != nil {
		logger.WithError(err).Warn("Error revoking role")
		utils.WriteError(w, http.StatusInternalServerError, "Error revoking role", nil)
	}

	logrus.WithField("userID", userID).Info("User revoked successfully : ", userRole)
	utils.WriteJson(w, http.StatusCreated, &ActionDeleted{
		Deleted: true,
	})

}

func (api *UserRoleAPI) GetRoleList(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user_role.go -> GetRoleList")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])

	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	userRoles, err := api.DB.GetRolesByUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("Error getting user roles")
		utils.WriteError(w, http.StatusInternalServerError, "Error getting user roles", nil)
		return
	}
	if userRoles == nil {
		userRoles = make([]*model.UserRole, 0)
	}
	logrus.WithField("userRoles", userRoles).Info("User roles fetched successfully")
	utils.WriteJson(w, http.StatusOK, &userRoles)

}

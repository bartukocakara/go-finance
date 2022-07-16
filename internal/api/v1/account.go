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

type AccountAPI struct {
	DB database.Database
}

func SetAccountAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &AccountAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{UserID}/accounts", api.Create, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/users/{UserID}/accounts", api.List, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/users/{UserID}/accounts/{AccountID}", api.Get, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodPatch, "/users/{UserID}/accounts/{AccountID}", api.Update, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodDelete, "/users/{UserID}/accounts/{AccountID}", api.Delete, auth.Admin, auth.MemberIsTarget),
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{UserID}/accounts
// Permission - Admin, MemberIsTarget
func (api *AccountAPI) Create(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "account.go -> Create()")
	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"Principal": principal,
		"UserID":    userID,
	})
	//Decode body
	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		logger.WithError(err).Warn("Could not decode account parameters")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode parameters", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	account.UserID = &userID

	if err := account.Verify(); err != nil {
		logger.WithError(err).Warn("Error while creating account")
		utils.WriteError(w, http.StatusInternalServerError, "Error while creating account", nil)
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateAccount(ctx, &account); err != nil {
		logger.WithError(err).Warn("Error while creating account")
		utils.WriteError(w, http.StatusInternalServerError, "Error while creating account", nil)
		return
	}

	logger.WithField("AccountID", account.ID).Info("Success: Account created")
	utils.WriteJson(w, http.StatusCreated, &account)
}

// GET - /users/{UserID}/accounts
// Permissions - Admin, MemberISTarget
func (api *AccountAPI) List(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "account.go -> List()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"Principal": principal,
		"UserID":    userID,
	})

	ctx := r.Context()
	accounts, err := api.DB.ListAccountsByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("Error getting account list")
		utils.WriteError(w, http.StatusInternalServerError, "Error while getting account list", nil)
		return
	}

	logger.Info("Success : Account listed")

	utils.WriteJson(w, http.StatusOK, accounts)
}

// GET - /users/{UserID}/accounts/{AccountID}
// Permissions - Admin, MemberIsTarget
func (api *AccountAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "account.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	accountID := model.AccountID(vars["AccountID"])
	principal := auth.GetPrincipal(r)

	logger.WithFields(logrus.Fields{
		"UserID":    userID,
		"AccountID": accountID,
		"Principal": principal,
	})

	ctx := r.Context()
	account, err := api.DB.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Warn("Error while getting account")
		utils.WriteError(w, http.StatusBadRequest, "Error while getting account", nil)
		return
	}

	logger.Info("Success: Error get")
	utils.WriteJson(w, http.StatusOK, &account)
}

// PATCH - /users/{UserID}/accounts/{AccountID}
// Permissions - Admin, MemberIsTarget
func (api *AccountAPI) Update(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "account.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	accountID := model.AccountID(vars["AccountID"])
	principal := auth.GetPrincipal(r)
	logger = logger.WithFields(logrus.Fields{
		"Principal": principal,
		"UserID":    userID,
		"AccountID": accountID,
	})

	var accountRequest model.Account
	if err := json.NewDecoder(r.Body).Decode(&accountRequest); err != nil {
		logger.WithError(err).Warn("Could not decode account parameters")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode account update parameters", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	ctx := r.Context()
	account, err := api.DB.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Warn("Error getting account")
		utils.WriteError(w, http.StatusInternalServerError, "Error getting acocunt", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	account = account.SetAll(&accountRequest)
	if err := api.DB.UpdateAccountByID(ctx, account); err != nil {
		logger.WithError(err).Warn("Error while updating account")
		utils.WriteError(w, http.StatusInternalServerError, "Error while updating account", nil)
		return
	}
	logger.Info("Account updated")
	utils.WriteJson(w, http.StatusOK, &ActionUpdated{
		Updated: true,
	})
}

// DELETE - /users/{UserID}/accounts/{AccountID}
// Permissions - Admin, MemberIsTarget
func (api *AccountAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "account.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	accountID := model.AccountID(vars["AccountID"])
	principal := auth.GetPrincipal(r)
	logger.WithFields(logrus.Fields{
		"UserID":    userID,
		"AccountID": accountID,
		"Principal": principal,
	})

	ctx := r.Context()
	ok, err := api.DB.DeleteAccountByID(ctx, accountID)
	if !ok && err != nil {
		logger.WithError(err).Warn("Error while deleting account")
		utils.WriteError(w, http.StatusInternalServerError, "Error deleting account", nil)
		return
	}

	logger.Info("Success: Account deleted")
	utils.WriteJson(w, http.StatusOK, &ActionDeleted{
		Deleted: true,
	})
}

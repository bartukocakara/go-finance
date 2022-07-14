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

type MerchantAPI struct {
	DB database.Database
}

func SetMerchantAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := MerchantAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{UserID}/merchants", api.Create, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/users/{UserID}/merchants", api.List, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodPatch, "/users/{UserID}/merchants/{MerchantID}", api.Update, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodDelete, "/users/{UserID}/merchants/{MerchantID}", api.Delete, auth.Admin, auth.MemberIsTarget),
	}
	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{userID}/merchants
// Permission - MemberIsTarget
func (api *MerchantAPI) Create(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "merchant.go -> Ceate()")
	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":    userID,
		"principal": principal,
	})

	var merchant model.Merchant
	if err := json.NewDecoder(r.Body).Decode(&merchant); err != nil {
		logger.WithError(err).Warn("Could not decode parameters from request")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode parameters", nil)
		return
	}

	merchant.UserID = &userID
	ctx := r.Context()

	if err := merchant.Verify(); err != nil {
		logger.WithError(err).Warn("Not all fields found")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := api.DB.CreateMerchant(ctx, &merchant); err == database.MerchantNameExists {
		logger.WithError(err).Warn("Merchant name already exists")
		utils.WriteError(w, http.StatusInternalServerError, "Merchant name already exists", nil)
		return
	} else if err != nil {
		logger.WithError(err).Warn("Error creating merchant")
		utils.WriteError(w, http.StatusInternalServerError, "Error Creatin Merchant", nil)
		return
	}

	logger.WithField("MerchantID", merchant.ID).Info("Merchant Created")
	utils.WriteJson(w, http.StatusCreated, &merchant)
}

// GET - /users/{userID}/merchants
// Permission - MemberIsTarget
func (api *MerchantAPI) List(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "merchant.go -> List()")
	vars := mux.Vars(r)

	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()
	merchants, err := api.DB.ListMerchantsByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("Error getting merchants")
		utils.WriteError(w, http.StatusInternalServerError, "Error while getting merchants", nil)
		return
	}
	logger.Info("Merchants returned")
	utils.WriteJson(w, http.StatusOK, &merchants)
}

// PATCH - users/{userID}/merchants/{merchantID}
// Permission - MemberIsTarget
func (api *MerchantAPI) Update(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "merchant.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["MserID"])
	merchantID := model.MerchantID(vars["MerchantID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"merchantID": merchantID,
	})

	// Decode paramters
	var merchantRequest model.Merchant
	if err := json.NewDecoder(r.Body).Decode(&merchantRequest); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("merchantID", merchantID)

	ctx := r.Context()

	merchant, err := api.DB.GetMerchantByID(ctx, merchantID)
	if err != nil {
		logger.WithError(err).Warn("error getting merchant")
		utils.WriteError(w, http.StatusConflict, "error getting merchant", nil)
		return
	}

	if merchantRequest.Name != nil || len(*merchantRequest.Name) != 0 {
		merchant.Name = merchantRequest.Name
	}

	// update merchant
	if err := api.DB.UpdateMerchant(ctx, merchant); err != nil {
		logger.WithError(err).Warn("Error while updating merchant")
		utils.WriteError(w, http.StatusInternalServerError, "Error updating merchant", nil)
		return
	}

	// log update status
	logger.Info("Merchant Updated")
	// http response user
	utils.WriteJson(w, http.StatusOK, &merchant)

}

func (api *MerchantAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "merchant.go -> Delete()")

	vars := mux.Vars(r)

	merchantID := model.MerchantID(vars["MerchantID"])

	principal := auth.GetPrincipal(r)
	logger = logger.WithFields(logrus.Fields{
		"merchantID": merchantID,
		"principal":  principal,
	})
	ctx := r.Context()
	// check merchant exists
	_, err := api.DB.GetMerchantByID(ctx, merchantID)
	if err != nil {
		logger.WithError(err).Warn("Error merchant not found")
		utils.WriteError(w, http.StatusBadRequest, "Error merchant not found", nil)
		return
	}

	// delete merchant by Id
	if err := api.DB.DeleteMerchantByID(ctx, &merchantID); err != nil {
		logger.WithError(err).Warn("Error could not delete merchant")
		utils.WriteError(w, http.StatusInternalServerError, "Error merchant could not deleted", nil)
		return
	}

	logger.Info("Merchant deleted")
	utils.WriteJson(w, http.StatusOK, &ActionDeleted{
		Deleted: true,
	})

}

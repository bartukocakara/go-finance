package v1

import (
	"encoding/json"
	"financial-app/internal/api/auth"
	"financial-app/internal/api/utils"
	"financial-app/internal/database"
	"financial-app/internal/model"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CategoryAPI struct {
	DB database.Database
}

func SetCategoryAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := CategoryAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{UserID}/categories", api.Create, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/users/{UserID}/categories", api.List, auth.MemberIsTarget),
		NewAPI(http.MethodPatch, "/users/{UserID}/categories/{CategoryID}", api.Update, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/users/{UserID}/categories/{CategoryID}", api.Get, auth.MemberIsTarget),
		NewAPI(http.MethodDelete, "/users/{UserID}/categories/{CategoryID}", api.Delete, auth.MemberIsTarget),
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{UserID}/categories
// Permission - MemberIsTarget
func (api *CategoryAPI) Create(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "category.go -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)
	logger = logger.WithFields(logrus.Fields{
		"UserID":    userID,
		"principal": principal,
	})

	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		logger.WithError(err).Error("Could not decode category params")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode params", map[string]string{
			"error": err.Error(),
		})
		return
	}

	category.UserID = &userID
	if err := category.Verify(); err != nil {
		logger.WithError(err).Warn("Not all fields found for category")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields found for category", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateCategory(ctx, &category); err != nil {
		logger.WithError(err).Warn("Error creating category")
		utils.WriteError(w, http.StatusInternalServerError, "Error creating category", nil)
		return
	}

	logger.WithField("catgoryID", category.ID).Info("Category Created ")
	utils.WriteJson(w, http.StatusCreated, &category)
}

// GET - /users/{UserID}/categories
// Permission - MemberIsTarget
func (api *CategoryAPI) List(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "categories.go -> List()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger = logrus.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	categories, err := api.DB.ListCategoriesByUserID(ctx, userID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting categories %s", userID)
		logger.WithError(err).Error(errorMessage)
		utils.WriteError(w, http.StatusBadRequest, errorMessage, nil)
		return
	}

	logger.Info("Success : Listing categories", userID)
	utils.WriteJson(w, http.StatusOK, &categories)
}

// PATCH - /users/{UserID}/categories/{CategoryID}
// Permission - MemberIsTarget
func (api *CategoryAPI) Update(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{})

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	categoryID := model.CategoryID(vars["CategoryID"])

	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"categoryID": categoryID,
	})

	var categoryRequest model.Category
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		logger.WithError(err).Error("Could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("CategoryID", categoryID)
	ctx := r.Context()

	category, err := api.DB.GetCategoryByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Warn("Error getting category")
		utils.WriteError(w, http.StatusInternalServerError, "Error getting category", nil)
		return
	}
	category.ParentID = category.SetParentID(categoryRequest.ParentID)
	category.Name = categoryRequest.SetName(categoryRequest.Name)

	if err := api.DB.UpdateCategory(ctx, category); err != nil {
		logger.WithError(err).Error("Error updating category")
		utils.WriteError(w, http.StatusBadRequest, "Error updating category", nil)
		return
	}
	logger.Info("Success : Category updated")
	utils.WriteJson(w, http.StatusOK, &ActionUpdated{
		Updated: true,
	})
}

// GET - /users/{UserID}/categories/{CategoryID}
// Permission - MemberIsTarget
func (api *CategoryAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "category.go -> Get()")
	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	categoryID := model.CategoryID(vars["CategoryID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":     userID,
		"principal":  principal,
		"categoryID": categoryID,
	})

	ctx := r.Context()

	category, err := api.DB.GetCategoryByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Warn("Error while getting category")
		utils.WriteError(w, http.StatusConflict, "Error while getting category", nil)
		return
	}

	logger.Info("Category returned")
	utils.WriteJson(w, http.StatusOK, &category)
}

func (api *CategoryAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "category.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	categoryID := model.CategoryID(vars["CategoryID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"UserID":     userID,
		"CategoryID": categoryID,
		"principal":  principal,
	})

	ctx := r.Context()

	ok, err := api.DB.DeleteCategoryByID(ctx, categoryID)
	if !ok && err != nil {
		logger.WithError(err).Warn("Error deleting category")
		utils.WriteError(w, http.StatusConflict, "Error deleting category", nil)
		return
	}

	logger.Info("Success: Category deleted")
	utils.WriteJson(w, http.StatusOK, &ActionDeleted{
		Deleted: true,
	})
}

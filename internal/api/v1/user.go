package v1

import (
	"context"
	"encoding/json"
	"financial-app/internal/api/auth"
	"financial-app/internal/api/utils"
	"financial-app/internal/database"
	"financial-app/internal/model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserAPI struct {
	DB database.Database // Will represent all database interface
}

func SetUserAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &UserAPI{
		DB: db,
	}
	apis := []API{
		// -----------USER----------------------------
		NewAPI(http.MethodPost, "/users", api.Create, auth.Any),
		NewAPI(http.MethodGet, "/users", api.List, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodPost, "/login", api.Login, auth.Any),
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, api.Func).Methods(api.Method)
	}
}

type UserParameters struct {
	model.User
	model.SessionData

	Password string `json:"password"`
}

func (api *UserAPI) Create(w http.ResponseWriter, r *http.Request) {
	// Show function name to track error faster
	logger := logrus.WithField("func", "user.go -> Create()")

	// Load parameters
	var userParameters UserParameters
	if err := json.NewDecoder(r.Body).Decode(&userParameters); err != nil {
		logger.WithError(err).Warn("Could not decode request user parameters")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode parameters", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"email": *userParameters.Email,
	})

	if err := userParameters.User.Verify(); err != nil {
		logger.WithError(err).Warn("Not all fields found", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	hashed, err := model.HashPassword(userParameters.Password)
	if err != nil {
		logger.WithError(err).Warn("Could not hash password")
		utils.WriteError(w, http.StatusInternalServerError, "Coul not hash password", nil)
		return
	}

	newUser := &model.User{
		Email:    userParameters.Email,
		Password: &hashed,
	}

	ctx := r.Context()
	if err := api.DB.CreateUser(ctx, newUser); err == database.ErrUserExists {
		logger.WithError(err).Warn("User already exists")
		utils.WriteError(w, http.StatusConflict, "User already exists", nil)
		return
	} else if err != nil {
		logger.WithError(err).Warn("Error creating user")
		utils.WriteError(w, http.StatusConflict, "Error creating user", nil)
		return
	}

	createdUser, err := api.DB.GetUserByID(ctx, newUser.ID)
	if err != nil {
		logger.WithError(err).Warn("Error creating user")
		utils.WriteError(w, http.StatusConflict, "Error creating user", nil)
		return
	}

	logger.WithField("UserID", createdUser.ID).Info("User created")
	api.WriteTokenResponse(ctx, w, http.StatusCreated, createdUser, &userParameters.SessionData, true)

}

type TokenResponse struct {
	Token *auth.Token `json:"token,omitempty"` // this will insert all tokens struct fields
	User  *model.User `json:"user,omitempty"`
}

func (api *UserAPI) WriteTokenResponse(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	user *model.User,
	sessionData *model.SessionData,
	cookie bool) {

	token, err := auth.IssueToken(model.Principal{UserID: user.ID})
	if err != nil && token == nil {
		logrus.WithError(err).Warn("Error issuing token")
		utils.WriteError(w, http.StatusConflict, "Error issuing token", nil)
		return
	}

	session := model.Session{
		UserID:       user.ID,
		DeviceID:     sessionData.DeviceID,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.RefreshTokenExpiresAt,
	}

	if err := api.DB.SaveRefreshToken(ctx, session); err != nil {
		logrus.WithError(err).Warn("Error issuing token")
		utils.WriteError(w, http.StatusConflict, "Error issuing token", nil)
		return
	}

	tokenResponse := TokenResponse{
		Token: token,
		User:  user,
	}
	if cookie {
		// later
	}

	utils.WriteJson(w, status, tokenResponse)
}

func (api *UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go -> Login()")

	var credentials model.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		logger.WithError(err).Warn("Could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "Could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := credentials.SessionData.Verify(); err != nil {
		logger.WithError(err).Warn("Not all fields found")
		utils.WriteError(w, http.StatusBadRequest, "Not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}
	logger = logger.WithFields(logrus.Fields{
		"email": credentials.Email,
	})

	ctx := r.Context()
	user, err := api.DB.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		logger.WithError(err).Warn("Error login in")
		utils.WriteError(w, http.StatusConflict, "Invalid password", nil)
		return
	}

	logger.WithField("userID", user.ID).Info("User Login")
	api.WriteTokenResponse(ctx, w, http.StatusOK, user, &credentials.SessionData, true)
}

func (api *UserAPI) List(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("Func", "user.go ->List()")

	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"principal": principal,
	})

	ctx := r.Context()

	users, err := api.DB.ListUsers(ctx)
	if err != nil {
		logger.WithError(err).Warn("Error getting users")
		utils.WriteError(w, http.StatusConflict, "Error getting users", nil)
		return
	}

	logger.Info("Users returned")
	utils.WriteJson(w, http.StatusOK, &users)
}

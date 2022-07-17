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

type TransactionAPI struct {
	DB database.Database
}

func SetTransactionAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &TransactionAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{UserID}/transactions", api.Create, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/transactions", api.ListAll, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/transactions/{TransactionID}", api.Get, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/accounts/{AccountID}/transactions", api.ListByAccount, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/categories/{CategoryID}/transactions", api.ListByCategory, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodGet, "/users/{UserID}/transactions", api.ListByUser, auth.Admin, auth.MemberIsTarget),
		NewAPI(http.MethodPatch, "/users/{UserID}/transactions/{TransactionID}", api.Update, auth.Admin),
		NewAPI(http.MethodDelete, "/users/{UserID}/transactions/{TransactionID}", api.Delete, auth.Admin),
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{UserID}/transactions
// Permissions - Admin, MemberIsTarget

func (api *TransactionAPI) Create(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transaction.go  -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)

	logger.WithFields(logrus.Fields{
		"Principal": principal,
		"UserID":    userID,
	})

	// decode params
	var transactionRequest model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		logger.WithError(err).Warn("Error could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "Error could not decode parameters", map[string]string{
			"Error": err.Error(),
		})
		return
	}
	transaction := transactionRequest.SetAll(&userID, &transactionRequest)
	if err := transaction.Verify(); err != nil {
		logger.WithError(err).Warn("Could not verified params")
		utils.WriteError(w, http.StatusBadRequest, "Error could not verify params", map[string]string{
			"Error": err.Error(),
		})
		return
	}
	ctx := r.Context()
	if err := api.DB.CreateTransaction(ctx, transaction); err != nil {
		logger.WithError(err).Warn("Error: Could not inserted transaction data")
		utils.WriteError(w, http.StatusInternalServerError, "Error: Could not inserted data", nil)
		return
	}

	logger.Info("Success: Transaction inserted")
	utils.WriteJson(w, http.StatusOK, transaction)
}

// GET - /transactions
// Permissions - Admin,memberIsTarget
func (api *TransactionAPI) ListAll(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transactions.go -> ListAll()")

	principal := auth.GetPrincipal(r)
	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Warn("Error: Invalid from parameters")
		utils.WriteError(w, http.StatusConflict, "Error : Invalid from parameters", nil)
		return
	}
	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Warn("Error: Invalid to parameter")
		utils.WriteError(w, http.StatusConflict, "Error : Invalid to parameters", nil)
		return
	}
	logger.WithFields(logrus.Fields{
		"Principal": principal,
	})
	ctx := r.Context()
	transactions, err := api.DB.ListAllTransactions(ctx, from, to)
	if err != nil {
		logger.WithError(err).Warn("Error: While getting transactions all")
		utils.WriteError(w, http.StatusConflict, "Error : While getting transactions", nil)
		return
	}

	logger.Info("Success: List All Categories")
	if transactions == nil {
		transactions = make([]*model.Transaction, 0)
	}
	utils.WriteJson(w, http.StatusOK, &transactions)

}

// GET - /transactions/{TransactionID}
// Permissions - admin,memberIsTarget
func (api *TransactionAPI) Get(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transaction.go -> Get()")

	vars := mux.Vars(r)
	transactionID := model.TransactionID(vars["TransactionID"])
	principal := auth.GetPrincipal(r)
	logger.WithFields(logrus.Fields{
		"Principal":     principal,
		"TransactionID": transactionID,
	})
	ctx := r.Context()
	transaction, err := api.DB.GetTransactionByID(ctx, transactionID)
	if err != nil {
		logger.WithError(err).Warn("Error while getting transaction")
		utils.WriteError(w, http.StatusConflict, "Error while getting transaction", nil)
		return
	}

	logger.Info("Success: Get Transaction")
	utils.WriteJson(w, http.StatusOK, &transaction)
}

// GET - /categories/{CategoryID}/transactions
// Permissions - admin,memberIsTarget
func (api *TransactionAPI) ListByCategory(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transactions.go -> ListByCategory")

	vars := mux.Vars(r)
	categoryID := model.CategoryID(vars["CategoryID"])
	principal := auth.GetPrincipal(r)

	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Warn("Error: Invalid from parameters")
		utils.WriteError(w, http.StatusConflict, "Error : Invalid from parameters", nil)
		return
	}
	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Warn("Error: Invalid to parameter")
		utils.WriteError(w, http.StatusConflict, "Error : Invalid to parameters", nil)
		return
	}

	logger.WithFields(logrus.Fields{
		"Principal":  principal,
		"CategoryID": categoryID,
		"From":       from,
		"To":         to,
	})

	ctx := r.Context()
	transactions, err := api.DB.ListTransactionByCategoryID(ctx, categoryID, from, to)
	if err != nil {
		logger.WithError(err).Warn("Error: While getting transaction by category")
		utils.WriteError(w, http.StatusConflict, "Error: while getting transaction by category", nil)
		return
	}

	logger.Info("Success: Transactions returned")
	if transactions == nil {
		transactions = make([]*model.Transaction, 0)
	}

	utils.WriteJson(w, http.StatusOK, &transactions)
}

// GET - /accounts/{AccountID}/transactions
// Methods - admin, memberIsTarget
func (api *TransactionAPI) ListByAccount(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transactions.go -> ListByAccount()")

	vars := mux.Vars(r)
	accountID := model.AccountID(vars["AccountID"])
	principal := auth.GetPrincipal(r)
	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Warn("Error: Invalid from parameters")
		utils.WriteError(w, http.StatusConflict, "Error : Invalid from parameters", nil)
		return
	}
	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Warn("Error: Invalid to parameter")
		utils.WriteError(w, http.StatusConflict, "Error : Invalid to parameters", nil)
		return
	}
	logger.WithFields(logrus.Fields{
		"Principal": principal,
		"AccountID": accountID,
		"From":      from,
		"To":        to,
	})
	ctx := r.Context()
	transactions, err := api.DB.ListTransactionByAccountID(ctx, accountID, from, to)
	if err != nil {
		logger.WithError(err).Warn("Error: While listing transactions by AccountID")
		utils.WriteError(w, http.StatusConflict, "Error :While listing transactions by AccountID", nil)
		return
	}

	logger.Info("Success: List transactions by AccountID")
	if transactions == nil {
		transactions = make([]*model.Transaction, 0)
	}
	utils.WriteJson(w, http.StatusOK, &transactions)
}

// GET /users/{UserID}/transactions
// Permissions - admin, memberIsTarget
func (api *TransactionAPI) ListByUser(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transactions.go -> ListTransactionByUserID()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	principal := auth.GetPrincipal(r)
	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Error("Error while parsing from query")
		utils.WriteError(w, http.StatusBadRequest, "Error while parsing from query", nil)
		return
	}
	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Error("Error while parsing to query")
		utils.WriteError(w, http.StatusBadRequest, "Error while parsing from query", nil)
		return
	}
	logger.WithFields(logrus.Fields{
		"Principal": principal,
		"UserID":    userID,
		"From":      from,
		"To":        to,
	})
	ctx := r.Context()
	transactions, err := api.DB.ListTransactionByUserID(ctx, userID, from, to)
	if err != nil {
		logger.WithError(err).Warn("Error while getting transaction by userID")
		utils.WriteError(w, http.StatusBadRequest, "Error while getting transaction by user ID", nil)
		return
	}
	logger.Info("Success: List transaction by user ID")
	utils.WriteJson(w, http.StatusOK, transactions)

}

// PATCH - /users/{UserID}/transactions/{TransactionID}
// Permissions - admin, memberIsTarget
func (api *TransactionAPI) Update(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transactions.go -> Update()")
	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	transactionID := model.TransactionID(vars["TransactionID"])
	principal := auth.GetPrincipal(r)
	logger.WithFields(logrus.Fields{
		"Principal":     principal,
		"TransactionID": transactionID,
		"UserID":        userID,
	})

	var transactionRequest model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		logger.WithError(err).Warn("Error : while parsing update transcation")
		utils.WriteError(w, http.StatusBadRequest, "Error : while parsing update transaction", map[string]string{
			"Error": err.Error(),
		})
		return
	}
	ctx := r.Context()
	transactionRequest.UserID = &userID
	if err := transactionRequest.Verify(); err != nil {
		logger.WithError(err).Warn("Error: Missing fields")
		utils.WriteError(w, http.StatusBadRequest, "Error: missing fields", nil)
		return
	}
	transaction, err := api.DB.GetTransactionByID(ctx, transactionID)
	if err != nil {
		logger.WithError(err).Warn("Error: Missing fields")
		utils.WriteError(w, http.StatusInternalServerError, "Error : transaction not found", nil)
		return
	}

	transaction = transaction.SetAll(&userID, &transactionRequest)
	if err := api.DB.UpdateTransaction(ctx, transaction); err != nil {
		logger.WithError(err).Warn("error updating transaction")
		utils.WriteError(w, http.StatusInternalServerError, "error updating transaction", nil)
		return
	}

	logger.Info("Success: Transaction updated")
	utils.WriteJson(w, http.StatusOK, &ActionUpdated{
		Updated: true,
	})

}

// DELETE - /users/{UserID}/transactions/{TransactionID}
// Permissions - admin, memberIsTarget
func (api *TransactionAPI) Delete(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "transaction.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["UserID"])
	transactionID := model.TransactionID(vars["TransactionID"])
	principal := auth.GetPrincipal(r)
	logger.WithFields(logrus.Fields{
		"Principal":     principal,
		"UserID":        userID,
		"TransactionID": transactionID,
	})

	ctx := r.Context()
	if err := api.DB.DeleteTransaction(ctx, transactionID); err != nil {
		logger.WithError(err).Warn("error deleting transaction")
		utils.WriteError(w, http.StatusConflict, "error deleting transaction", nil)
		return
	}

	logger.Info("Success: Transaction deleted")
	utils.WriteJson(w, http.StatusOK, &ActionDeleted{
		Deleted: true,
	})
}

package v1

import (
	"financial-app/internal/api/auth"
	"financial-app/internal/database"

	"github.com/gorilla/mux"
)

type MerchantAPI struct {
	DB database.Database
}

func SetMerchantAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {

}

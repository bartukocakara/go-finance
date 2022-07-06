package auth

import (
	"financial-app/internal/model"
	"net/http"
)

type principalContextKeyType struct{}

var principalContextKey principalContextKeyType

func GetPrincipal(r *http.Request) model.Principal {
	if principal, ok := r.Context().Value(principalContextKey).(model.Principal); ok {
		return principal
	}
	return model.NilPrincipal
}

package v1

import (
	"financial-app/internal/api/auth"
	"net/http"
)

type API struct {
	Path            string
	Method          string
	Func            http.HandlerFunc
	permissionTypes []auth.PermissionType
}

func NewAPI(method string, path string, handlerFunc http.HandlerFunc, permissionTypes ...auth.PermissionType) API {
	return API{
		Path:            path,
		Method:          method,
		Func:            handlerFunc,
		permissionTypes: permissionTypes,
	}
}

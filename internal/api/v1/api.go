package v1

import "net/http"

type API struct {
	Path   string
	Method string
	Func   http.HandlerFunc
}

func NewAPI(method string, path string, handlerFunc http.HandlerFunc) API {
	return API{
		Path:   path,
		Method: method,
		Func:   handlerFunc,
	}
}

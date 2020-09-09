package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Name(route.Name).
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"WrapList",
		"GET",
		"/pasteles",
		WrapList,
	},
	Route{
		"WrapShow",
		"GET",
		"/pastel/{id}",
		WrapShow,
	},
	Route{
		"WrapAdd",
		"POST",
		"/pastel",
		WrapAdd,
	},
	Route{
		"WrapUpdate",
		"PUT",
		"/pastel/{id}",
		WrapUpdate,
	},
	Route{
		"WrapRemove",
		"DELETE",
		"/pastel/{id}",
		WrapRemove,
	},
}

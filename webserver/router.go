package webserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

type HttpRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerHandleFunc
}

type HttpRoutes []HttpRoute

func newRouter(env *Environment, httpRoutes HttpRoutes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// setup normal http routes
	for _, route := range httpRoutes {
		var handler http.Handler
		handler = Handler{Env: env, Handle: route.HandlerFunc}
		handler = Logger(handler, route.Name)

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}

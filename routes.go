package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"ArticleIndex",
		"GET",
		"/v0/articles",
		ArticleIndex,
	},
	Route{
		"FilterIndex",
		"GET",
		"/v0/filters",
		FilterIndex,
	},
	Route{
		"TodoIndex",
		"GET",
		"/v0/todos",
		TodoIndex,
	},
	Route{
		"TodoCreate",
		"POST",
		"/v0/todos",
		TodoCreate,
	},
	Route{
		"TodoShow",
		"GET",
		"/v0/todos/{todoId}",
		TodoShow,
	},
}

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
		"ArticleSearch",
		"POST",
		"/v0/articles/search",
		ArticleSearch,
	},
	Route{
		"FilterAll",
		"GET",
		"/v0/filters",
		FilterAll,
	},
	Route{
		"FilterById",
		"GET",
		"/v0/filter/{FilterId:[0-9]+}",
		FilterById,
	},
	Route{
		"Imgage",
		"GET",
		"/v0/Image/{format}p/{Path:[a-zA-Z0-9/_]+}/{FileName:[a-zA-Z0-9_]+.(JPG|jpg)}",
		Image,
	},
}

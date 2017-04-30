package main

import (
	"github.com/geknuepft/backend/image"
	"github.com/geknuepft/backend/webserver"
)

var httpRoutes = webserver.HttpRoutes{
	webserver.HttpRoute{
		"Index",
		"GET",
		"/",
		webserver.HandleIndex,
	},
	/*
	HttpRoute{
		"ArticleIndex",
		"GET",
		"/v0/articles",
		ArticleIndex,
	},
	HttpRoute{
		"ArticleDetailById",
		"GET",
		"/v0/article/{ArticleId:[0-9]+}",
		ArticleDetailById,
	},
	HttpRoute{
		"ArticleSearch",
		"POST",
		"/v0/articles",
		ArticleSearch,
	},
	HttpRoute{
		"FilterAll",
		"GET",
		"/v0/filters",
		FilterAll,
	},
	HttpRoute{
		"FilterById",
		"GET",
		"/v0/filter/{FilterId:[0-9]+}",
		FilterById,
	},
	*/
	webserver.HttpRoute{
		"Imgage",
		"GET",
		"/v0/Image/{format}p/{Path:[a-zA-Z0-9/_]+}/{FileName:[a-zA-Z0-9_]+.(?:JPG|jpg)}",
		image.ImageHandleGet,
	},
}

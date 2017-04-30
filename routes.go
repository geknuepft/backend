package main

import (
	"github.com/geknuepft/backend/image"
	"github.com/geknuepft/backend/webserver"
	"github.com/geknuepft/backend/articles"
)

var httpRoutes = webserver.HttpRoutes{
	webserver.HttpRoute{
		"Index",
		"GET",
		"/",
		webserver.HandleIndex,
	},
	webserver.HttpRoute{
		"ArticleIndex",
		"GET",
		"/v0/articles",
		articles.ArticleIndex,
	},
	webserver.HttpRoute{
		"ArticleDetailById",
		"GET",
		"/v0/articles/{ArticleId:[0-9]+}",
		articles.ArticleDetailById,
	},
	webserver.HttpRoute{
		"ArticleSearch",
		"POST",
		"/v0/articles",
		articles.ArticleSearch,
	},
	webserver.HttpRoute{
		"FilterAll",
		"GET",
		"/v0/filters",
		articles.FilterAll,
	},
	webserver.HttpRoute{
		"FilterById",
		"GET",
		"/v0/filter/{FilterId:[0-9]+}",
		articles.FilterById,
	},
	webserver.HttpRoute{
		"Imgage",
		"GET",
		"/v0/Image/{format}p/{Path:[a-zA-Z0-9/_]+}/{FileName:[a-zA-Z0-9_]+.(?:JPG|jpg)}",
		image.ImageHandleGet,
	},
}

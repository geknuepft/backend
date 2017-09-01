package main

import (
	"github.com/geknuepft/backend/image"
	"github.com/geknuepft/backend/webserver"
	"github.com/geknuepft/backend/api"
)

var httpRoutes = webserver.HttpRoutes{
	webserver.HttpRoute{
		"Index",
		"GET",
		"/",
		webserver.HandleIndex,
	},
	webserver.HttpRoute{
		"Imgage",
		"GET",
		"/v0/Image/{format}p/{Path:[a-zA-Z0-9/_]+}/{FileName:[a-zA-Z0-9_]+.(?:JPG|jpg)}",
		image.ImageHandleGet,
	},
	webserver.HttpRoute{
		"InstnaceGet",
		"GET",
		"/v0/Instance",
		api.InstanceGet,
	},
}

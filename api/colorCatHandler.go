package api

import (
	"github.com/geknuepft/backend/webserver"
	"net/http"
)

func ColorCatGet(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {

	colorCats, err := getColorCats()

	if err != nil {
		return webserver.StatusError{500, err}
	}

	return webserver.WriteJson(w, colorCats)
}

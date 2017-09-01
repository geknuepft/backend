package api

import (
	"github.com/geknuepft/backend/webserver"
	"net/http"
	"strconv"
)

func InstanceGet(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {

	// parse get parameters
	_, err := strconv.Atoi(r.URL.Query().Get("length_mm"))
	if err != nil {

	}

	instances := GetInstances()
	return webserver.WriteJson(w, instances)
}
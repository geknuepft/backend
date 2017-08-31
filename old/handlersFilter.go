package old

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/geknuepft/backend/webserver"
)

func FilterAll(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	filters := GetFilters()
	return webserver.WriteJson(w, filters)
}

func FilterById(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var filterId int
	var err error
	if filterId, err = strconv.Atoi(vars["FilterId"]); err != nil {
		return webserver.StatusError{400, err}
	}

	filter, err := GetFilterById(filterId)
	if err != nil {
		return webserver.StatusError{404, err}
	}

	return webserver.WriteJson(w, filter)
}

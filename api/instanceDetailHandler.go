package api

import (
	"github.com/geknuepft/backend/webserver"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

func InstanceDetailGet(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	var instanceId int
	var err error
	if instanceId, err = strconv.Atoi(vars["InstanceId"]); err != nil {
		return webserver.StatusError{400, err}
	}

	instanceDetail, err := GetInstanceDetailById(instanceId)

	if err != nil {
		return webserver.StatusError{404, err}
	}

	return webserver.WriteJson(w, instanceDetail)
}

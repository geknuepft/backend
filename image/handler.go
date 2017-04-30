package image

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"log"
	"github.com/geknuepft/backend/webserver"
)

func ImageHandleGet(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var format int
	var err error
	if format, err = strconv.Atoi(vars["format"]); err != nil {
		panic(err)
	}

	path := vars["Path"] + "/" + vars["FileName"]

	oupImg, err := ImageGet(format, path)
	if err != nil {
		return webserver.StatusError{404, err}
	}

	webserver.WriteJpegHeaders(w, http.StatusOK)
	err = ImageWrite(w, oupImg)
	if err != nil {
		return webserver.StatusError{500, err}
	}

	log.Printf("serve image: format=%v, path=%v", format, path)

	return nil
}

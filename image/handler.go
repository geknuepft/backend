package image

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/geknuepft/backend/webserver"
	"os"
	"errors"
	"fmt"
	"time"
)

func writeImageHeaders(w http.ResponseWriter, lastModified time.Time) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", MaxAge))
	SetLastModified(w, lastModified)
}

func ImageHandleGet(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var format int
	var err error
	if format, err = strconv.Atoi(vars["format"]); err != nil {
		panic(err)
	}

	filePath := ImageInputDir + vars["Path"] + "/" + vars["FileName"]

	fileInfo, err := getFileNameAndStat(filePath)
	if err != nil {
		return webserver.StatusError{http.StatusNotFound, err}
	}

	// check if StatusNotModified can be answered
	lastModified := GetLastModified(fileInfo)
	if CheckIfModifiedSince(r, lastModified) {
		return webserver.StatusError{http.StatusNotModified, errors.New("not modified")}
	}

	oupImg, err := ImageGet(format, filePath)
	if err != nil {
		return webserver.StatusError{http.StatusInternalServerError, err}
	}

	writeImageHeaders(w, lastModified)
	err = ImageWrite(w, oupImg)
	if err != nil {
		return webserver.StatusError{500, err}
	}

	return nil
}

func getFileNameAndStat(filePath string) (fileInfo os.FileInfo, err error) {
	if fileInfo, err = os.Stat(filePath); err == nil {
		return
	}

	return nil,
		errors.New(fmt.Sprintf("cannot find filePath=%v", filePath))
}

package webserver

import (
	"net/http"
	"encoding/json"
	"log"
)

func WriteJson(w http.ResponseWriter, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return StatusError{500, err}
	}

	WriteJsonHeaders(w)
	w.Write(b)

	return nil
}

func WriteJsonError(w http.ResponseWriter, v interface{}, errorCode int ) error {
	w.WriteHeader(errorCode)
	WriteJsonHeaders(w);

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("WriteJsonError encoding problem: %v", err)
	}
	return nil
}
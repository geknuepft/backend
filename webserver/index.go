package webserver

import (
	"net/http"
	"fmt"
)

func HandleIndex(env *Environment, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Welcome to the geknuepf backend server!\n")
	return nil
}

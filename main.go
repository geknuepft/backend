package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// http://thenewstack.io/make-a-restful-json-api-go/

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

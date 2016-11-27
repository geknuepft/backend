package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	router := NewRouter()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("while parsing the PORT env variable:")
		log.Fatal(err.Error())
		os.Exit(2)
	}

	log.Printf("[geknuepft-backend] listening on port %v", port)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

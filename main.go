package main

import (
	"log"
	"os"
	"strconv"
	"github.com/geknuepft/backend/webserver"
)

func main() {
	// get configuration
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("while parsing the PORT env variable:")
		log.Fatal(err.Error())
		os.Exit(2)
	}

	// start http server
	env := &webserver.Environment{}
	webserver.Run("", port, env, httpRoutes)
}

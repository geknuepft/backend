package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the geknuepf backend server!\n")
}

func ArticleIndex(w http.ResponseWriter, r *http.Request) {
	articles := GetArticles()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(articles, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func ArticleSearch(w http.ResponseWriter, r *http.Request) {
	// read request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// parse reuqest body
	var filterValues FilterValues

	if err := json.Unmarshal(body, &filterValues); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// check against json schema
	schemaLoader := gojsonschema.NewStringLoader(`
{
    "title": "v0/articles/search Schema",
    "type": "array",
    "items": {
        "type": "object",
        "properties": {
            "filter_id": {
                "type": "integer"
            },
            "values": {
                "type": "array",
                "items": {
                    "type": "integer"
                }
            }
        },
        "required": [
            "filter_id",
            "values"
        ]
    }
}
    `)
	documentLoader := gojsonschema.NewStringLoader(string(body))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if !result.Valid() {
		response := gojsonschemaErrorToSring(result.Errors())

		w.WriteHeader(422) // Unprocessable Entity
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}

		return
	}

	// get data
	articles := GetArticlesByFilterValues(filterValues)

	// write resposne
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(articles, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func gojsonschemaErrorToSring(inp []gojsonschema.ResultError) (oup []string) {

	oup = make([]string, len(inp))

	for i, error := range inp {
		oup[i] = error.String()
	}

	return
}

func FilterAll(w http.ResponseWriter, r *http.Request) {
	filters := GetFilters()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(filters, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func FilterById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	var filterId int
	var err error
	if filterId, err = strconv.Atoi(vars["FilterId"]); err != nil {
		panic(err)
	}

	filters := GetFilterById(filterId)
	if len(filters) > 0 {
		w.WriteHeader(http.StatusOK)

		b, err := json.MarshalIndent(filters[0], "", "    ")
		if err != nil {
			panic(err)
		}
		w.Write(b)
	} else {
		jsonWriteNotFound(w)
	}
}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func jsonWriteNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)

	ret := jsonErr{Code: http.StatusNotFound, Text: "Not Found"}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

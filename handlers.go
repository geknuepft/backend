package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the geknuepf backend server!\n")
}

func writeJsonHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)
}

func writeJpegHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week
	w.WriteHeader(status)
}

func ArticleDetailById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var articleId int
	var err error
	if articleId, err = strconv.Atoi(vars["ArticleId"]); err != nil {
		panic(err)
	}

	articleDetail, err := GetArticleDetailById(articleId)

	if err == nil {
		writeJsonHeaders(w, http.StatusOK)

		b, err := json.MarshalIndent(articleDetail, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Write(b)
	} else {
		jsonWriteNotFound(w)
	}
}

func ArticleSearch(w http.ResponseWriter, r *http.Request) {
	// parse get parameters
	lengthMm, err := strconv.Atoi(r.URL.Query().Get("length_mm"))
	if err != nil {
		lengthMm = 140
	}

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
	articles := GetArticlesByFilterValues(filterValues, lengthMm)

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

	writeJsonHeaders(w, http.StatusOK)

	b, err := json.MarshalIndent(filters, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func FilterById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var filterId int
	var err error
	if filterId, err = strconv.Atoi(vars["FilterId"]); err != nil {
		panic(err)
	}

	filter, err := GetFilterById(filterId)

	if err == nil {
		writeJsonHeaders(w, http.StatusOK)

		b, err := json.MarshalIndent(filter, "", "    ")
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
	jsonWriteError(w, "Object Not Found")
}

func jsonWriteError(w http.ResponseWriter, text string) {
	writeJsonHeaders(w, http.StatusNotFound)

	ret := jsonErr{Code: http.StatusNotFound, Text: text}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func Image(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var format int
	var err error
	if format, err = strconv.Atoi(vars["format"]); err != nil {
		panic(err)
	}

	path := vars["Path"] + "/" + vars["FileName"]

	oupImg, err := ImageGet(format, path)

	if err == nil {
		writeJpegHeaders(w, http.StatusOK)
		err = ImageWrite(w, oupImg)
	}

	if err == nil {
		log.Printf("serve image: format=%v, path=%v", format, path)
	} else {
		jsonWriteError(w, err.Error())
	}
}

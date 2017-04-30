package articles

import (
	"github.com/xeipuuv/gojsonschema"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"github.com/geknuepft/backend/webserver"
)

func ArticleIndex(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	// parse get parameters
	lengthMm, err := strconv.Atoi(r.URL.Query().Get("length_mm"))
	if err != nil {
		lengthMm = 140
	}

	articles := GetArticlesByFilterValues(FilterValues{}, lengthMm)
	return webserver.WriteJson(w, articles)
}

func ArticleDetailById(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var articleId int
	var err error
	if articleId, err = strconv.Atoi(vars["ArticleId"]); err != nil {
		panic(err)
	}

	articleDetail, err := GetArticleDetailById(articleId)
	return webserver.WriteJson(w, articleDetail)
}

func ArticleSearch(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	// read request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return webserver.StatusError{400, err}
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		return webserver.StatusError{400, err}
	}

	// parse reuqest body
	var filterValues FilterValues

	if err := json.Unmarshal(body, &filterValues); err != nil {
		return webserver.WriteJsonError(w, err, http.StatusBadRequest)
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
		return webserver.StatusError{500, err}
	}

	if !result.Valid() {
		response := gojsonschemaErrorToSring(result.Errors())

		// Unprocessable Entity
		return webserver.WriteJsonError(w, response, 422)
	}

	// get data
	articles := GetArticlesByFilterValues(filterValues, 180)

	// write resposne
	webserver.WriteJson(w, articles)
	return nil
}

func gojsonschemaErrorToSring(inp []gojsonschema.ResultError) (oup []string) {
	oup = make([]string, len(inp))

	for i, err := range inp {
		oup[i] = err.String()
	}

	return
}

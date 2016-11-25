package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the geknuepf backend server!\n")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId int
	var err error
	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
		panic(err)
	}
	todo := RepoFindTodo(todoId)
	if todo.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos
*/
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
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

		b, err := json.MarshalIndent(filters, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Write(b)
	} else {
		jsonWriteNotFound(w)
	}
}

func jsonWriteNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)

	ret := jsonErr{Code: http.StatusNotFound, Text: "Not Found"}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

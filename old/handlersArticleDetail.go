package old

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/geknuepft/backend/webserver"
)

func ArticleDetailById(env *webserver.Environment, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var articleId int
	var err error
	if articleId, err = strconv.Atoi(vars["ArticleId"]); err != nil {
		panic(err)
	}

	articleDetail, err := GetArticleDetailById(articleId)
	if err != nil {
		return webserver.StatusError{404, err}
	}

	return webserver.WriteJson(w, articleDetail)
}

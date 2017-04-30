package webserver

import "net/http"

func WriteJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Model", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
}

func WriteJpegHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week
}

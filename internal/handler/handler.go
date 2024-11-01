package handler

import (
	"net/http"
)

func HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}
}

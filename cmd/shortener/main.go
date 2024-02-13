package main

import (
	"github.com/ilramm/gourlshortener/internal/handlers"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.HandlePost(w, r)
		case http.MethodGet:
			handlers.HandleGet(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}

package main

import (
	"fmt"
	"net/http"
)

var urlID string

func handleGet(wr http.ResponseWriter, rr *http.Request) {

	originalGetURL := rr.URL.String()
	urlID = originalGetURL[1:]
	host := rr.Host

	wr.Header().Add("Location", "https://practicum.yandex.ru/")
	wr.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Fprintf(wr, "%s\n%s", urlID, host)

}

func handlePost(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	originalPostURL := r.Header.Get("URL")

	if contentType != "text/plain" {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	fmt.Fprintf(w, "Received data: %s, %s", contentType, originalPostURL)
	w.WriteHeader(http.StatusCreated)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePost(w, r)
		case http.MethodGet:
			handleGet(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

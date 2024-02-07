package main

import (
	"io"
	"math/rand"
	"net/http"
	"sync"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var urlKeys sync.Map

func generateRandomString(length int) string {

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)

}

func handlePost(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)

	if string(body) != "" {
		randomString := generateRandomString(5)
		theURL := string(body)
		urlKeys.Store(randomString, theURL)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("URL", theURL)
		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + randomString))
	}

}

func handleGet(wr http.ResponseWriter, rr *http.Request) {

	urlID := rr.URL.String()[1:]

	if originalURL, ok := urlKeys.Load(urlID); ok {

		str := originalURL.(string)
		wr.Header().Set("Location", str)
		wr.WriteHeader(307)

	} else {

		wr.Header().Set("Location", "URL not found")
		wr.WriteHeader(400)

	}
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
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}

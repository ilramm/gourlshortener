package main

import (
	"math/rand"
	"net/http"
	"strings"
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

	theURL := r.Header.Get("URL")
	host := r.Host
	randomString := generateRandomString(5)

	urlKeys.Store(randomString, theURL)
	finalURL := host + "/" + randomString
	w.Header().Set("URL", finalURL)
	w.WriteHeader(http.StatusCreated)

}

func handleGet(wr http.ResponseWriter, rr *http.Request) {

	urlID := rr.URL.String()[1:]

	if originalURL, ok := urlKeys.Load(urlID); ok {

		str := originalURL.(string)

		if strings.Contains(str, "https://") {
			wr.Header().Set("Location", str)
			wr.WriteHeader(307)
		} else {
			newStr := "https://" + str
			wr.Header().Set("Location", newStr)
			wr.WriteHeader(307)
		}

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

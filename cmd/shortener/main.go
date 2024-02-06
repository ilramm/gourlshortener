package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var urlKeys map[string]string

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

	urlKeys = make(map[string]string, 10)
	urlKeys[randomString] = theURL

	finalURL := host + "/" + randomString
	w.Header().Set("URL", finalURL)
	fmt.Fprintf(w, "%s", urlKeys[randomString])
	w.WriteHeader(http.StatusCreated)

}

func handleGet(wr http.ResponseWriter, rr *http.Request) {
	originalGetURL := rr.URL.String()
	urlID := originalGetURL[1:]
	originalURL := urlKeys[urlID]

	// Debugging statements
	fmt.Fprintf(wr, "originalGetURL: %s, urlID: %s, originalURL: %s\n", originalGetURL, urlID, originalURL)

	// Check if originalURL is empty
	if originalURL == "" {
		http.Error(wr, "Original URL not found", http.StatusNotFound)
		return
	}

	// Set headers and perform redirect
	wr.Header().Add("Location", originalURL)
	wr.WriteHeader(http.StatusTemporaryRedirect)

	// Optional: You can include a message in the response body
	fmt.Fprintf(wr, "Redirecting to: %s", originalURL)
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

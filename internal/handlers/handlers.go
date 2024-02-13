package handlers

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

func HandlePost(w http.ResponseWriter, r *http.Request) {

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

func HandleGet(wr http.ResponseWriter, rr *http.Request) {

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

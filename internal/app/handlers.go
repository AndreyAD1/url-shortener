package app

import (
	"io"
	"net/http"
	"net/url"
	"path"
)

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createShortURLHandler(w, r)
		return
	}
	if r.Method == http.MethodGet {
		getFullURLHandler(w, r)
		return
	}
	errMsg := "Only GET and POST requests are allowed"
	http.Error(w, errMsg, http.StatusMethodNotAllowed)
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	incomingURL, err := url.Parse(string(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL := GetShortURL(*incomingURL)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func getFullURLHandler(w http.ResponseWriter, r *http.Request) {
	path, urlID := path.Split(r.URL.Path)
	if path != "/" || urlID == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	fullURL, err := GetFullURL(urlID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

package handlers

import (
	"io"
	"net/http"
	"net/url"
	"path"

	srv "github.com/AndreyAD1/url-shortener/internal/app/service"
)

func CreateShortURLHandler(service srv.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		incomingURL, err := url.ParseRequestURI(string(requestBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortURL, err := service.GetShortURL(*incomingURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}
}

func GetFullURLHandler(service srv.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path, urlID := path.Split(r.URL.Path)
		if path != "/" || urlID == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		fullURL, err := service.GetFullURL(urlID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

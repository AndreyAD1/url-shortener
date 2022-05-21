package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/gorilla/mux"
)

func (h HandlerContainer) CreateShortURLHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		incomingURL := string(requestBody)
		if _, err := url.ParseRequestURI(incomingURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortURL, err := h.URLService.CreateShortURL(incomingURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}
}

func (h HandlerContainer) GetFullURLHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urlID := mux.Vars(r)["id"]
		fullURL, err := h.URLService.GetFullURL(urlID)
		if errors.Is(err, service.ErrorNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h HandlerContainer) CreateShortURLApiHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var requestInfo CreateShortURLRequest
		if err := json.Unmarshal(requestBody, &requestInfo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := url.ParseRequestURI(requestInfo.URL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortURL, err := h.URLService.CreateShortURL(requestInfo.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response{Result: shortURL})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

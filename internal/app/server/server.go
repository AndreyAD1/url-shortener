package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func NewServer(address string) *http.Server {
	return &http.Server{Addr: address, Handler: GetHandler()}
}

func GetHandler() http.Handler {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	router := mux.NewRouter()
	router.HandleFunc(
		"/",
		handlers.CreateShortURLHandler(URLService),
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/{id}",
		handlers.GetFullURLHandler(URLService),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/api/shorten_url",
		handlers.CreateShortURLApiHandler(URLService),
	)
	return router
}

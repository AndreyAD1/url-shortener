package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func NewServer(cfg config.StartupConfig) *http.Server {
	return &http.Server{Addr: cfg.ServerAddress, Handler: GetHandler(cfg)}
}

func GetHandler(cfg config.StartupConfig) http.Handler {
	db, err := storage.NewStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	URLService := service.Service{
		Storage:        db,
		BaseURL:        cfg.BaseURL,
		ShortURLLength: cfg.ShortURLLength,
	}
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
		"/api/shorten",
		handlers.CreateShortURLApiHandler(URLService),
	)
	return router
}

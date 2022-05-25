package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
	"github.com/AndreyAD1/url-shortener/internal/app/middlewares"
	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func NewServer(cfg config.StartupConfig) (*http.Server, error) {
	handler, err := GetHandler(cfg)
	if err != nil {
		return nil, err
	}
	return &http.Server{Addr: cfg.ServerAddress, Handler: handler}, nil
}

func GetHandler(cfg config.StartupConfig) (http.Handler, error) {
	db, err := storage.NewStorage(cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}
	URLService := service.Service{
		Storage:        db,
		BaseURL:        cfg.BaseURL,
		ShortURLLength: cfg.ShortURLLength,
	}
	handlers := handlers.HandlerContainer{URLService: URLService}
	router := mux.NewRouter()
	router.HandleFunc(
		"/",
		handlers.CreateShortURLHandler(),
	).Methods(http.MethodPost)
	router.HandleFunc(
		"/{id}",
		handlers.GetFullURLHandler(),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/api/shorten",
		handlers.CreateShortURLApiHandler(),
	)
	router.Use(middlewares.Auth)
	router.Use(middlewares.DecompressGzipRequest)
	router.Use(middlewares.CompressResponseToGzip)
	return router, nil
}

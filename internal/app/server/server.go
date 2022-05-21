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

func NewServer(cfg config.StartupConfig) *http.Server {
	return &http.Server{Addr: cfg.ServerAddress, Handler: GetHandler(cfg)}
}

func GetHandler(cfg config.StartupConfig) http.Handler {
	db := storage.NewStorage(cfg.FileStoragePath)
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
	router.Use(middlewares.DecompressGzipRequest)
	router.Use(middlewares.CompressResponseToGzip)
	return router
}

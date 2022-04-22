package server

import (
	"net/http"

	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func NewServer(address string) *http.Server {
	return &http.Server{Addr: address, Handler: getHandler()}
}

func getHandler() http.Handler {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	handler := http.NewServeMux()
	handlerFunc := http.HandlerFunc(handlers.ShortURLHandler(URLService))
	handler.Handle("/", handlerFunc)
	return handler
}

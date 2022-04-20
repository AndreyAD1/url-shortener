package server

import (
	"net/http"

	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func NewServer(address string) *http.Server {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	handler := http.HandlerFunc(handlers.ShortURLHandler(URLService))
	http.HandleFunc("/", handler)
	return &http.Server{Addr: address}
}

package server

import (
	"net/http"
	
	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
)

func NewServer(address string) *http.Server {
	http.HandleFunc("/", handlers.ShortURLHandler)
	return &http.Server{Addr: address}
}

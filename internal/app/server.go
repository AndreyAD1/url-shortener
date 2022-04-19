package app

import "net/http"

func NewServer(address string) *http.Server {
	http.HandleFunc("/", ShortURLHandler)
	return &http.Server{Addr: address}
}

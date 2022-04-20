package main

import (
	"log"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/server"
)

func main() {
	srv := server.NewServer(config.ServerAddress)
	err := srv.ListenAndServe()
	log.Println(err)
}

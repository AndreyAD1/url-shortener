package main

import (
	"log"

	"github.com/AndreyAD1/url-shortener/internal/app"
)

func main() {
	server := app.NewServer(app.ServerAddress)
	err := server.ListenAndServe()
	log.Println(err)
}

package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/server"
)

func main() {
	serverAddress := flag.String("a", "", "a server address")
	baseURL := flag.String("b", "", "a shorten URL host")
	fileStoragePath := flag.String("f", "", "a path to a file storage")
	cfg := config.StartupConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	flag.Parse()
	if *serverAddress != "" {
		cfg.ServerAddress = *serverAddress
	}
	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if *fileStoragePath != "" {
		cfg.FileStoragePath = *fileStoragePath
	}
	srv := server.NewServer(cfg)
	err := srv.ListenAndServe()
	log.Println(err)
}

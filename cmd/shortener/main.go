package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/server"
)

func main() {
	cfg := config.StartupConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	srv := server.NewServer(cfg)
	err := srv.ListenAndServe()
	log.Println(err)
}

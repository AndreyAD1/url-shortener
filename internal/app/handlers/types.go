package handlers

import 	srv "github.com/AndreyAD1/url-shortener/internal/app/service"


type HandlerContainer struct {
	URLService srv.Service
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type Response struct {
	Result interface{} `json:"result"`
}
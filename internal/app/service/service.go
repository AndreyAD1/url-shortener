package service

import (
	"math/rand"
	"time"

	s "github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Service struct {
	Storage        s.Repository
	BaseURL        string
	ShortURLLength int
}

func GetRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s Service) GetShortURL(url string) (string, error) {
	randomString := GetRandomString(s.ShortURLLength)
	err := s.Storage.WriteURL(randomString, url)
	if err != nil {
		return "", err
	}
	shortURL := s.BaseURL + "/" + randomString
	return shortURL, nil
}

func (s Service) GetFullURL(urlID string) (string, error) {
	fullURL, err := s.Storage.GetURL(urlID)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}

package service

import (
	"fmt"
	"math/rand"
	u "net/url"
	"time"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	s "github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Service struct {
	Storage s.Storage
}

func GetRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s Service) GetShortURL(url u.URL) (string, error) {
	randomString := GetRandomString(config.ShortURLLength)
	err := s.Storage.WriteURL(randomString, url)
	if err != nil {
		return "", err
	}
	shortURL := "http://" + config.ServerAddress + "/" + randomString
	return shortURL, nil
}

func (s Service) GetFullURL(urlID string) (string, error) {
	fullURL, err := s.Storage.GetURL(urlID)
	if err != nil {
		return "", fmt.Errorf("URL not found")
	}
	return fullURL.String(), nil
}

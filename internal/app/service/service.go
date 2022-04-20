package service

import (
	"fmt"
	"math/rand"
	u "net/url"
	"time"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetShortURL(url u.URL) string {
	randomString := GetRandomString(config.ShortURLLength)
	shortURL := "http://" + config.ServerAddress + "/" + randomString
	storage.URLStorage[randomString] = url
	return shortURL
}

func GetFullURL(urlID string) (string, error) {
	fullURL, ok := storage.URLStorage[urlID]
	if !ok {
		return "", fmt.Errorf("URL not found")
	}
	return fullURL.String(), nil
}

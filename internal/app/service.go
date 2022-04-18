package app

import (
	"fmt"
	"math/rand"
	u "net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetShortURL(url u.URL) string {
	randomString := getRandomString(ShortURLLength)
	shortURL := "http://" + ServerAddress + "/" + randomString
	URLStorage[randomString] = url
	return shortURL
}

func GetFullURL(urlID string) (string, error) {
	fullURL, ok := URLStorage[urlID]
	if !ok {
		return "", fmt.Errorf("URL not found")
	}
	return fullURL.String(), nil
}

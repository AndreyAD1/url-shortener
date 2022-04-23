package storage

import (
	"fmt"
	"net/url"
)

type Repository interface {
	WriteURL(urlID string, fullURL url.URL) error
	GetURL(urlID string) (url.URL, error)
}

type Storage map[string]url.URL

func (s Storage) WriteURL(urlID string, fullURL url.URL) error {
	s[urlID] = fullURL
	return nil
}

func (s Storage) GetURL(urlID string) (url.URL, error) {
	fullURL, ok := s[urlID]
	if !ok {
		return fullURL, fmt.Errorf("no URL was found")
	}
	return fullURL, nil
}

func NewStorage() Storage {
	return Storage(make(map[string]url.URL))
}

var URLStorage = make(map[string]url.URL)

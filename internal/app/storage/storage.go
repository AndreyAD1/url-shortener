package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Repository interface {
	WriteURL(urlID string, fullURL string) error
	GetURL(urlID string) (string, error)
}

type FileStorage struct {
	filename string
}

type MemoryStorage map[string]string

type URLInfo struct {
	ID  string
	URL string
}

func (s MemoryStorage) WriteURL(urlID string, fullURL string) error {
	s[urlID] = fullURL
	return nil
}

func (s MemoryStorage) GetURL(urlID string) (string, error) {
	fullURL, ok := s[urlID]
	if !ok {
		return fullURL, fmt.Errorf("no URL was found")
	}
	return fullURL, nil
}

func (s FileStorage) WriteURL(urlID string, fullURL string) error {
	fileFlag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile(s.filename, fileFlag, 0777)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	URLItem := URLInfo{ID: urlID, URL: fullURL}
	if err := encoder.Encode(URLItem); err != nil {
		log.Printf("storage encoder error: %v", err)
		return err
	}
	return nil
}

func (s FileStorage) GetURL(urlID string) (string, error) {
	fileFlag := os.O_RDONLY | os.O_CREATE
	file, err := os.OpenFile(s.filename, fileFlag, 0777)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(file)
	for {
		URLItem := &URLInfo{}
		err := decoder.Decode(&URLItem)
		if err != nil {
			log.Printf("storage decoder error: %v", err)
			break
		}
		if URLItem.ID == urlID {
			return URLItem.URL, nil
		}
	}
	return "", fmt.Errorf("no URL was found")
}

func NewStorage(storageFile string) Repository {
	if storageFile == "" {
		return MemoryStorage(make(map[string]string))
	}
	return FileStorage{filename: storageFile}
}

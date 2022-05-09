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
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
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
	URLItem := URLInfo{ID: urlID, URL: fullURL}
	if err := s.encoder.Encode(URLItem); err != nil {
		return err
	}
	return nil
}

func (s FileStorage) GetURL(urlID string) (string, error) {
	storageContent := make(map[string]string)
	for {
		var URLItem URLInfo
		err := s.decoder.Decode(&URLItem)
		if err != nil {
			log.Printf("storage decoder error: %v", err)
			break
		}
		storageContent[URLItem.ID] = URLItem.URL
	}
	fullURL, ok := storageContent[urlID]
	if !ok {
		return fullURL, fmt.Errorf("no URL was found")
	}
	return fullURL, nil
}

func NewStorage(storageFile string) (Repository, error) {
	if storageFile == "" {
		return MemoryStorage(make(map[string]string)), nil
	}
	fileFlag := os.O_RDWR | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile(storageFile, fileFlag, 0777)
	if err != nil {
		return FileStorage{}, err
	}
	fileStorage := FileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}
	return fileStorage, err
}

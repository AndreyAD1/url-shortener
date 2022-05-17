package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Repository interface {
	WriteURL(urlID string, fullURL string) error
	GetURL(urlID string) (string, error)
}

type FileStorage struct {
	filename string
	sync.RWMutex
}

type MemoryStorage struct {
	storage map[string]string
	sync.RWMutex
}
type URLInfo struct {
	ID  string
	URL string
}

func (s *MemoryStorage) WriteURL(urlID string, fullURL string) error {
	s.Lock()
	defer s.Unlock()
	s.storage[urlID] = fullURL
	return nil
}

func (s *MemoryStorage) GetURL(urlID string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	fullURL, ok := s.storage[urlID]
	if !ok {
		return fullURL, fmt.Errorf("no URL was found")
	}
	return fullURL, nil
}

func (s *FileStorage) WriteURL(urlID string, fullURL string) error {
	s.Lock()
	defer s.Unlock()
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

func (s *FileStorage) GetURL(urlID string) (string, error) {
	s.RLock()
	defer s.RUnlock()
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
	var mu sync.RWMutex
	if storageFile == "" {
		storage := MemoryStorage{make(map[string]string), mu}
		return &storage
	}
	storage := FileStorage{storageFile, mu}
	return &storage
}

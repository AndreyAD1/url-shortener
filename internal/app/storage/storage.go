package storage

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Repository interface {
	WriteURL(urlID string, fullURL string) error
	GetURL(urlID string) (*string, error)
}

type FileStorage struct {
	filename string
	storage  map[string]string
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

func (s *MemoryStorage) GetURL(urlID string) (*string, error) {
	s.RLock()
	defer s.RUnlock()
	fullURL, ok := s.storage[urlID]
	if !ok {
		return nil, nil
	}
	return &fullURL, nil
}

func (s *FileStorage) WriteURL(urlID string, fullURL string) error {
	s.Lock()
	defer s.Unlock()
	fileFlag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile(s.filename, fileFlag, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	URLItem := URLInfo{ID: urlID, URL: fullURL}
	if err := encoder.Encode(URLItem); err != nil {
		log.Printf("storage encoder error: %v", err)
		return err
	}
	s.storage[URLItem.ID] = URLItem.URL
	return nil
}

func (s *FileStorage) GetURL(urlID string) (*string, error) {
	s.RLock()
	defer s.RUnlock()
	fullURL, ok := s.storage[urlID]
	if !ok {
		return nil, nil
	}
	return &fullURL, nil
}

func (s *FileStorage) readFile() error {
	s.RLock()
	defer s.RUnlock()
	fileFlag := os.O_RDONLY | os.O_CREATE
	file, err := os.OpenFile(s.filename, fileFlag, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for {
		URLItem := &URLInfo{}
		err := decoder.Decode(&URLItem)
		if err != nil {
			log.Printf("storage decoder error: %v", err)
			break
		}
		s.storage[URLItem.ID] = URLItem.URL
	}
	return nil
}

func NewStorage(storageFile string) Repository {
	if storageFile == "" {
		storage := MemoryStorage{storage: make(map[string]string)}
		return &storage
	}
	storage := FileStorage{
		filename: storageFile,
		storage:  make(map[string]string),
	}
	storage.readFile()
	return &storage
}

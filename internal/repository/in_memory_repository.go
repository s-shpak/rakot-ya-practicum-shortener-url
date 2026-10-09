package repository

import (
	"errors"
	"fmt"
	"sync"
)

var ErrNotFound = errors.New("resource not found")

type URLRepository interface {
	Save(key, url string) error
	Find(key string) (string, error)
}

// InMemoryRepository реализует URLRepository для хранения данных в оперативной памяти.
type InMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		store: make(map[string]string),
	}
}

func (r *InMemoryRepository) Save(key, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[key] = url
	return nil
}

func (r *InMemoryRepository) Find(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.store[key]
	if !exists {
		return "", fmt.Errorf("key not found %s: %w", key, ErrNotFound)
	}
	return url, nil
}

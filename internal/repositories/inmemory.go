package repositories

import (
	"context"
	"sync"
)

// Структура хранилища InMemory
type inMemoryStorage struct {
	memory map[string]string
	mutex  sync.RWMutex
}

// Фунцкия создания хранилища в InMemory
func NewInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		memory: make(map[string]string),
	}
}

// Метод получения оригинальной ссылки из хранилища InMemory
func (m *inMemoryStorage) Get(_ context.Context, shortLink string) (*string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	originalLink, ok := m.memory[shortLink]
	if !ok {
		return nil, nil
	}
	return &originalLink, nil
}

// Метод записи в хранилище InMemory
func (m *inMemoryStorage) Save(_ context.Context, shortLink, originaLink string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.memory[shortLink] = originaLink
	return nil
}

// Метод проверки существования ссылки в InMemory
func (m *inMemoryStorage) CheckExistance(_ context.Context, shortLink string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.memory[shortLink]
	if !ok {
		return false, nil
	}
	return true, nil
}

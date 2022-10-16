package storage

import (
	"sync"
)

// Store is the thread safe in memory key value store.
type Store struct {
	sync.RWMutex
	values map[string]interface{}
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{
		values: make(map[string]interface{}),
	}
}

// Load returns the value for the specified key.
func (s *Store) Load(key string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	result, ok := s.values[key]
	return result, ok
}

// Remove removed the given key.
func (s *Store) Remove(key string) {
	delete(s.values, key)
}

// Exist checks if the given key exists.
func (s *Store) Exist(key string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.values[key]
	return ok
}

// Save persists the give key/vale combination.
func (s *Store) Save(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()
	s.values[key] = value
	return nil
}

// LoadAll is returning all the key/value from the store.
// It returns a map of keys to values.
func (s *Store) LoadAll() (map[string]interface{}, error) {
	s.RLock()
	defer s.RUnlock()
	copyValues := make(map[string]interface{}, len(s.values))
	for k, v := range s.values {
		copyValues[k] = v
	}
	return copyValues, nil
}

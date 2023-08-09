package container

import "sync"

type SyncMap[K comparable, V any] struct {
	sync.RWMutex
	data map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		data: make(map[K]V),
	}
}

func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	s.RLock()
	val, ok := s.data[key]
	s.RUnlock()
	return val, ok
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

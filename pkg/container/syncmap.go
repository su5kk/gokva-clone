package container

import "sync"

type SyncMap[K comparable, V any] struct {
	mx   sync.RWMutex
	data map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		data: make(map[K]V),
	}
}

func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	s.mx.RLock()
	val, ok := s.data[key]
	s.mx.RUnlock()
	return val, ok
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.mx.Lock()
	s.data[key] = value
	s.mx.Unlock()
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.mx.Lock()
	delete(s.data, key)
	s.mx.Unlock()
}

package lazywriter

import "sync"

// Cache -
type Cache struct {
	store map[string]*LazyObject
	m     sync.Mutex
}

// NewStore -
func NewStore() *Cache {
	return &Cache{
		store: make(map[string]*LazyObject),
	}
}

// Add -
func (s *Cache) Add(key string, value *LazyObject) {
	s.m.Lock()
	s.store[key] = value
	s.m.Unlock()
}

// MustGet -
func (s *Cache) MustGet(key string) *LazyObject {
	s.m.Lock()
	v := s.store[key]
	s.m.Unlock()
	return v
}

// Get -
func (s *Cache) Get(key string) (*LazyObject, bool) {
	s.m.Lock()
	v, ok := s.store[key]
	s.m.Unlock()
	return v, ok
}

// Delete -
func (s *Cache) Delete(key string) {
	s.m.Lock()
	delete(s.store, key)
	s.m.Unlock()
}

// Keys -
func (s *Cache) Keys() []string {
	s.m.Lock()
	keys, i := make([]string, len(s.store)), 0
	for id := range s.store {
		keys[i] = id
		i++
	}
	s.m.Unlock()
	return keys[:]
}

// // Iterate -
// func (s *Cache) Iterate(fn func(key string, value *LazyObject) error) {
// 	s.m.Lock()
// 	for id := range s.store {
// 		if err := fn(id, s.MustGet(id)); err != nil {
// 			// log error
// 		}
// 	}
// }

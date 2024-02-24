package db

import "sync"

type Store struct {
	index  map[int]bool
	cached []int
	mu     sync.RWMutex
}

func (s *Store) Insert(value int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.index[value] = true
}

func (s *Store) Read() []int {

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.cached) == len(s.index) {
		return s.cached
	}

	result := make([]int, len(s.index))
	i := 0
	for k := range s.index {
		result[i] = k
		i++
	}

	s.cached = result
	return result
}

func (s *Store) Merge(data []int) {

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, value := range data {
		s.index[value] = true
	}
}

func NewStore() *Store {
	index := map[int]bool{}
	cached := []int{}
	return &Store{index: index, cached: cached}
}

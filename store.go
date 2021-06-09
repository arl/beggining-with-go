package shorten

import (
	"math/rand"
	"sync"
	"time"
)

// Store is a simple key value store.
type Store interface {
	// Add saves a new value in the store and returns the unique key.
	Add(value string) uint64

	// Load loads the value associated with a key (or false if it doesn't exist).
	Load(key uint64) (value string, ok bool)

	// ForEach calls f for each key-value pair.
	ForEach(f StoreForEachFunc)
}

// ForEachFunc is called with each key-value pairs in a Store, returns false to
// stop iteration.
type StoreForEachFunc func(key uint64, value string) bool

// MemoryStore is a concurrent-safe in-memory key-value store.
type MemoryStore struct {
	mu  sync.Mutex
	m   map[uint64]string
	rng *rand.Rand
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		m:   make(map[uint64]string),
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Add saves a new value in the store and returns the unique key.
func (s *MemoryStore) Add(value string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Loop until we find a key that is not used.
	var key uint64
	for {
		key = s.rng.Uint64()
		if _, ok := s.m[key]; !ok {
			break
		}
	}

	s.m[key] = value
	return key
}

// Load loads the value associated with a key (or false if it doesn't exist).
func (s *MemoryStore) Load(key uint64) (value string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok = s.m[key]
	return value, ok
}

// ForEach calls f for each key-value pair.
func (s *MemoryStore) ForEach(f StoreForEachFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range s.m {
		if !f(k, v) {
			break
		}
	}
}

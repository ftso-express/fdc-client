package storage

import "sync"

type cyclicItem[T any] struct {
	key   uint64
	value T
}

// Cyclic is a limited size storage. Keys are nonnegative integers. Item with key n is stored to n (mod size) together with the key.
type Cyclic[T any] struct {
	values []*cyclicItem[T]
	mu     *sync.RWMutex
}

// Size is the size of cyclic storage.
func (s Cyclic[T]) Size() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return uint64(len(s.values))
}

// Store stores value with key to key (mod size).
func (s Cyclic[T]) Store(key uint64, value T) {
	keyMod := key % s.Size()

	storedItem := &cyclicItem[T]{key: key, value: value}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[keyMod] = storedItem
}

// Get retrieves element from key (mod size) if the stored element has key key.
func (s Cyclic[T]) Get(key uint64) (T, bool) {
	var k T
	keyMod := key % s.Size()

	s.mu.RLock()
	defer s.mu.RUnlock()
	storedItem := s.values[keyMod]

	if storedItem == nil {
		return k, false
	}

	storedKey := storedItem.key

	if storedKey != key {
		return k, false
	}

	k = storedItem.value

	return k, true
}

// NewCyclic initializes a Cyclic storage with size.
func NewCyclic[T any](size uint64) Cyclic[T] {
	return Cyclic[T]{values: make([]*cyclicItem[T], size), mu: new(sync.RWMutex)}
}

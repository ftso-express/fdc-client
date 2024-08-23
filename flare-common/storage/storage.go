package storage

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type cyclicItem[T any, K constraints.Integer] struct {
	key   K
	value T
}

// Cyclic is a limited size storage. Keys are nonnegative integers. Item with key n is stored to n (mod size) together with the key.
type Cyclic[T any, K constraints.Integer] struct {
	values []*cyclicItem[T, K]
	mu     *sync.RWMutex
}

// Size is the size of cyclic storage.
func (s Cyclic[T, K]) Size() K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return K(len(s.values))
}

// Store stores value with key to key (mod size).
func (s Cyclic[T, K]) Store(key K, value T) {
	keyMod := key % s.Size()

	storedItem := &cyclicItem[T, K]{key: key, value: value}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[keyMod] = storedItem
}

// Get retrieves element from key (mod size) if the stored element has key key.
func (s Cyclic[T, K]) Get(key K) (T, bool) {
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
func NewCyclic[T any, K constraints.Integer](size int) Cyclic[T, K] {
	return Cyclic[T, K]{values: make([]*cyclicItem[T, K], size), mu: new(sync.RWMutex)}
}

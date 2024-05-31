package storage

type cyclicItem[T any] struct {
	key   uint64
	value T
}

// Cyclic is a limited size storage. Keys are nonnegative integers. Item with key n is stored to n (mod size) together with the key.
type Cyclic[T any] struct {
	values []*cyclicItem[T]
}

// Size is the size of cyclic storage.
func (s Cyclic[T]) Size() uint64 {
	return uint64(len(s.values))
}

// Store stores value with key to key (mod size).
func (s Cyclic[T]) Store(key uint64, value T) {
	keyMod := key % s.Size()

	storedItem := &cyclicItem[T]{key: key, value: value}

	s.values[keyMod] = storedItem
}

// Get retrieves element from key (mod size) if the stored element has key key.
func (s Cyclic[T]) Get(key uint64) (T, bool) {
	var k T
	keyMod := key % s.Size()

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
	return Cyclic[T]{values: make([]*cyclicItem[T], size)}
}

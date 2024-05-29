package storage

type cyclicItem[T any] struct {
	key   uint64
	value T
}

type Cyclic[T any] struct {
	items map[uint64]cyclicItem[T]
	size  uint64
}

func (s Cyclic[T]) Size() uint64 {

	return s.size
}

func (s Cyclic[T]) Store(key uint64, value T) {

	keyMod := key % s.size

	storedItem := cyclicItem[T]{key: key, value: value}

	s.items[keyMod] = storedItem
}

func (s Cyclic[T]) Get(key uint64) (T, bool) {

	var k T
	keyMod := key % s.size

	storedItem, exists := s.items[keyMod]

	if !exists {
		return k, false
	}

	storedKey := storedItem.key

	if storedKey != key {
		return k, false
	}

	k = storedItem.value

	return k, true

}

func NewCyclic[T any](size uint64) Cyclic[T] {

	items := map[uint64]cyclicItem[T]{}

	return Cyclic[T]{items: items, size: size}

}

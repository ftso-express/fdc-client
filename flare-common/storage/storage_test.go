package storage_test

import (
	"flare-common/storage"
	"testing"
)

func TestCyclicSimple(t *testing.T) {

	const size uint64 = 4

	stg := storage.NewCyclic[uint64](size)

	if stg.Size() != size {
		t.Error("wrong size")
	}

	getEmpty, exists := stg.Get(123)

	if exists || getEmpty != 0 {
		t.Error("wrong size")
	}

	stg.Store(1, 1)

	getOne, exists := stg.Get(1)

	if !exists || getOne != 1 {
		t.Error("not stored")
	}

	getOneWithFive, exists := stg.Get(5)

	if exists || getOneWithFive != 0 {
		t.Error("returned with wrong key")
	}

	//overwrite
	stg.Store(5, 5)

	getOneWithFive, exists = stg.Get(5)

	if !exists || getOneWithFive != 5 {
		t.Error("not overwritten")
	}

	getOne, exists = stg.Get(1)

	if exists || getOne != 0 {
		t.Error("not overwritten with larger")
	}

}

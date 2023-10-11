package storagemgr

import (
	"sync"
)

type MemoryStorage[T any] struct {
	mutex sync.Mutex
	data  T
}

func (f *MemoryStorage[T]) ReadFile(path string) (T, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.data, nil
}

func (f *MemoryStorage[T]) WriteFile(path string, data T) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.data = data

	return nil
}

func NewMemoryStorage[T any]() StorageMgr[T] {
	return &MemoryStorage[T]{mutex: sync.Mutex{}, data: *new(T)}
}

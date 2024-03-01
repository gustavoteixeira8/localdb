package localdb

import (
	"encoding/json"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type StorageMgr[T any] interface {
	ReadFile(path string) (T, error)
	WriteFile(path string, data T) error
}

// JSON Implementation --------------------------------------------

type JSONStorage[T any] struct {
	mutex sync.Mutex
}

func (f *JSONStorage[T]) ReadFile(path string) (T, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	tval := new(T)

	databs, err := os.ReadFile(path)
	if err != nil {
		return *tval, err
	}

	err = json.Unmarshal(databs, &tval)
	if err != nil {
		return *tval, err
	}

	return *tval, nil
}

func (f *JSONStorage[T]) WriteFile(path string, data T) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	bs, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bs, 0777)
	if err != nil {
		return err
	}

	return nil
}

func NewJSONStorage[T any]() StorageMgr[T] {
	return &JSONStorage[T]{mutex: sync.Mutex{}}
}

// YAML Implementation --------------------------------------------

type YAMLStorage[T any] struct {
	mutex sync.Mutex
}

func (f *YAMLStorage[T]) ReadFile(path string) (T, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	tval := new(T)
	databs, err := os.ReadFile(path)
	if err != nil {
		return *tval, err
	}

	err = yaml.Unmarshal(databs, &tval)
	if err != nil {
		return *tval, err
	}

	return *tval, nil
}

func (f *YAMLStorage[T]) WriteFile(path string, data T) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	bs, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bs, 0777)
	if err != nil {
		return err
	}

	return nil
}

func NewYAMLStorage[T any]() StorageMgr[T] {
	return &YAMLStorage[T]{mutex: sync.Mutex{}}
}

// Memory Implementation --------------------------------------------

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

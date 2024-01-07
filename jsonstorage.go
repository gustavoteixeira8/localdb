package localdb

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONStorage[T any] struct {
	mutex sync.Mutex
}

func (f JSONStorage[T]) ReadFile(path string) (T, error) {
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

func (f JSONStorage[T]) WriteFile(path string, data T) error {
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

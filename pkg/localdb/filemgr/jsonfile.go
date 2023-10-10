package filemgr

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONFile[T any] struct {
	mutex sync.Mutex
}

func (f JSONFile[T]) ReadFile(path string) (T, error) {
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

func (f JSONFile[T]) WriteFile(path string, data T) error {
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

func NewJSONFile[T any]() FileMgr[T] {
	return &JSONFile[T]{mutex: sync.Mutex{}}
}

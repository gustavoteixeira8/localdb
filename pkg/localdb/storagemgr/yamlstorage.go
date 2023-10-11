package storagemgr

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

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

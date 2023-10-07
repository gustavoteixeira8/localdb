package filemgr

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLFile[T any] struct {
}

func (f *YAMLFile[T]) ReadFile(path string) (T, error) {
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

func (f *YAMLFile[T]) WriteFile(path string, data T) error {
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

func NewYAMLFile[T any]() FileMgr[T] {
	return &YAMLFile[T]{}
}

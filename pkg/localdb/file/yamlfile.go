package file

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLFile[T any] struct {
}

func (f *YAMLFile[T]) ReadFile(path string) (*T, error) {
	databs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tval := new(T)

	err = yaml.Unmarshal(databs, &tval)
	if err != nil {
		return nil, err
	}

	return tval, nil
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

func NewYAMLFile[T any]() *YAMLFile[T] {
	return &YAMLFile[T]{}
}

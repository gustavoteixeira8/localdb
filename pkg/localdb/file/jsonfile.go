package file

import (
	"encoding/json"
	"os"
)

type JSONFile[T any] struct {
}

func (f JSONFile[T]) ReadFile(path string) (*T, error) {
	databs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tval := new(T)

	err = json.Unmarshal(databs, &tval)
	if err != nil {
		return nil, err
	}

	return tval, nil
}

func (f JSONFile[T]) WriteFile(path string, data T) error {
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

func NewJSONFile[T any]() *JSONFile[T] {
	return &JSONFile[T]{}
}
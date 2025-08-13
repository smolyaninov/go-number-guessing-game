package repo

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JSONRepository[T any] struct {
	filename string
}

func NewJSONRepository[T any](filename string) *JSONRepository[T] {
	return &JSONRepository[T]{filename: filename}
}

func (r *JSONRepository[T]) Load() (T, error) {
	var data T

	if _, err := os.Stat(r.filename); os.IsNotExist(err) {
		return data, nil
	}

	file, err := os.Open(r.filename)
	if err != nil {
		return data, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (r *JSONRepository[T]) Save(data T) error {
	if err := os.MkdirAll(filepath.Dir(r.filename), 0755); err != nil {
		return err
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tempFile := r.filename + ".tmp"
	if err := os.WriteFile(tempFile, content, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, r.filename)
}

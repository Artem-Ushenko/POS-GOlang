package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func newStore() *Store {
	return &Store{
		NextCustomerID: 1,
		NextProductID:  1,
		NextSaleID:     1,
	}
}

func loadStore(path string) (*Store, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newStore(), nil
		}
		return nil, err
	}
	defer file.Close()

	var store Store
	if err := json.NewDecoder(file).Decode(&store); err != nil {
		return nil, err
	}

	return &store, nil
}

func saveStore(path string, store *Store) error {
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, "data-*.json")
	if err != nil {
		return err
	}
	tempFileName := tempFile.Name()
	defer func() {
		_ = os.Remove(tempFileName)
	}()

	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(store); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	return os.Rename(tempFileName, path)
}

package main

import (
	"encoding/json"
	"errors"
	"os"
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
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(store)
}

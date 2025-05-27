package clientapp

import (
	"client/account"
	"encoding/json"
	"os"
)

func WriteDataToFile(data account.Data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func DestroyDataFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

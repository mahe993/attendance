/*
Package db provides functionality for reading and writing JSON data to and from files.

The db package includes methods for reading and writing JSON data to files. It utilizes the encoding/json package for marshaling and unmarshaling JSON.
*/
package db

import (
	"encoding/json"
	"errors"

	"os"

	"attendance.com/src/logger"
)

// The Read function reads JSON data from a specified file path and unmarshals it into the provided payload.
// It returns an error if the file cannot be read, is empty, or if unmarshalling fails.
func Read(filePath string, payload interface{}) error {
	bs, err := os.ReadFile(os.Getenv("APP_DB_PATH") + filePath)
	logger.Println(os.Getenv("APP_DB_PATH") + filePath)
	if err != nil {
		logger.Println(err)
		panic(errors.New("unable to read from file:" + filePath))
	}

	if len(bs) == 0 {
		return errors.New("empty document")
	}

	err = json.Unmarshal(bs, payload)
	if err != nil {
		logger.Println(err)
		panic(errors.New("unmarshalling JSON failed"))
	}

	return nil
}

// The Write function marshals the provided payload into JSON format and writes it to the specified file path.
// It returns an error if marshalling fails or if the file cannot be written.
func Write(payload interface{}, filePath string) error {
	bs, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		logger.Println(err)
		panic(errors.New("marshalling JSON failed"))
	}

	err = os.WriteFile(os.Getenv("APP_DB_PATH")+filePath, bs, 0644)
	if err != nil {
		logger.Println(err)
		panic(errors.New("unable to write to file:" + filePath))
	}

	return nil
}

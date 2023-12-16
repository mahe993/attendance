package db

import (
	"encoding/json"
	"errors"

	"os"

	"attendance.com/src/logger"
)

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

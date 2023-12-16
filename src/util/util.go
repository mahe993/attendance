/*
Package utils offers various utility functions that can be used throughout the application.
*/
package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"attendance.com/src/logger"
	"github.com/joho/godotenv"
)

func init() {
	logger.Println("Initializing envs...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logger.Println("Success!")
}

// Variables for env values
var (
	SUB_IP_1 = os.Getenv("VALID_IP_ADDR_1")
	SUB_IP_2 = os.Getenv("VALID_IP_ADDR_2")
)

// ValidateIP validates the IP address of the user to ensure it matches the config and returns a boolean indication and error.
func ValidateIP() (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Println(err)
		return false, err
	}

	for _, i := range interfaces {
		if i.Name == "en0" {
			byName, err := net.InterfaceByName(i.Name)
			if err != nil {
				logger.Println(err)
			}

			addresses, err := byName.Addrs()
			if err != nil {
				logger.Println(err)
			}

			for _, v := range addresses {
				addressSlice := strings.Split(v.String(), ".")

				// check if address slice == x.x.x.x
				if len(addressSlice) == 4 {
					// validate address
					if addressSlice[0] == SUB_IP_1 && addressSlice[1] == SUB_IP_2 {
						return true, nil
					}
					return false, nil
				}
			}
		}
	}

	return false, nil
}

func ReadCSV(file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// WriteCSV writes the given CSV data to the given file path.
// It returns a channel that can be used to wait for the write to complete.
// Saving a copy of csv uploads is low priority and thus does not panic on errors
func WriteCSV(filePath string, csvData [][]string) <-chan bool {
	done := make(chan bool)
	go func() {
		defer close(done)

		destFile, err := os.Create(filePath)
		if err != nil {
			logger.Println(fmt.Sprint("Error creating CSV file:", err))
			return
		}
		defer destFile.Close()

		// Write each row of the CSV data to the output file.
		csvWriter := csv.NewWriter(destFile)
		for _, line := range csvData {
			if err := csvWriter.Write(line); err != nil {
				logger.Println(fmt.Sprint("Error writing to CSV file:", err))
				return
			}
		}
		csvWriter.Flush()
	}()

	return done
}

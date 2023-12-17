/*
Package utils offers various utility functions that can be used throughout the application.
*/
package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"attendance.com/src/logger"
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
				validIPSlice := strings.Split(os.Getenv("VALID_IP_ADDR"), ".")

				// check if address slice == x.x.x.x
				if len(addressSlice) == 4 {
					octetRange1, err := strconv.Atoi(addressSlice[2])
					if err != nil {
						logger.Println(err)
						return false, err
					}
					octetRange2, err := strconv.Atoi(addressSlice[3][:3])
					if err != nil {
						logger.Println(err)
						return false, err
					}

					// validate address
					if addressSlice[0] == validIPSlice[0] &&
						addressSlice[1] == validIPSlice[1] &&
						octetRange1 >= 0 && octetRange1 <= 255 &&
						octetRange2 >= 0 && octetRange2 <= 255 {
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

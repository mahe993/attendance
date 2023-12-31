/*
Package utils provides various utility functions that can be used throughout the application.

The utils package includes functions for validating IP addresses, reading and writing CSV files, and other general-purpose utilities.
*/
package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"attendance.com/src/logger"
)

// ReadCSV reads CSV data from the given io.Reader.
// It returns a 2D slice of strings representing the CSV data and an error.
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

// ValidateClientIPHandler validates the IP address of the user to ensure it matches the config and returns a boolean indication and error.
// It returns true if the IP address is valid, false if it is not, and an error if one occurs.
// It is used to ensure that users are on the appropriate WIFI before checking in.
func ValidateClientIPHandler(r *http.Request) (bool, error) {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.Header.Get("CF-Connecting-IP")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	// unable to verify IP address
	if IPAddress == "" {
		return false, nil
	}

	// validate IP
	addressSlice := strings.Split(IPAddress, ".")
	validIPSlice := strings.Split(os.Getenv("VALID_IP_ADDR"), ".")

	// check if address slice == x.x.x.x
	if len(addressSlice) != 4 {
		return false, nil
	}
	octetRange1, err := strconv.Atoi(addressSlice[2])
	if err != nil {
		logger.Println(err)
		return false, err
	}
	octetRange2, err := strconv.Atoi(addressSlice[3])
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

	// IP not in range
	return false, nil

}

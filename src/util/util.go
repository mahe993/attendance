/*
Package utils offers various utility functions that can be used throughout the application.
*/
package utils

import (
	"encoding/csv"
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

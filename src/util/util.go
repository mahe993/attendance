/*
Package utils offers various utility functions that can be used throughout the application.
*/
package utils

import (
	"net"
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

				// check if address slice == x.x.x.x
				if len(addressSlice) == 4 {
					// validate address
					if addressSlice[0] == "192" && addressSlice[1] == "168" {
						return true, nil
					}
					return false, nil
				}
			}
		}
	}

	return false, nil
}

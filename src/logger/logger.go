/*
Package logger provides a custom logging utility with additional file name and line number information.

The logger package includes a Println function, which logs messages along with the file name and line number from where the logging function is called.
*/
package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

// The Println function logs a message along with the file name and line number from where it is called.
// It uses the runtime.Caller function to retrieve the caller's file name and line number.
func Println(msg interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	splitFile := strings.Split(file, "/")
	fileName := strings.Join(splitFile[6:], "/")

	// Log the message along with file name and line number
	log.Printf(":::%s:%d:::\n%v\n", fileName, line, msg)
	fmt.Println()
}

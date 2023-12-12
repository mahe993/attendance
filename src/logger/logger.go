package logger

import (
	"log"
	"runtime"
	"strings"
)

func Println(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	splitFile := strings.Split(file, "/")
	fileName := strings.Join(splitFile[6:], "/")

	// Log the message along with file name and line number
	log.Printf("%s:%d - %s\n", fileName, line, msg)
}

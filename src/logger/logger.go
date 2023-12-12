package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Println(msg interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	splitFile := strings.Split(file, "/")
	fileName := strings.Join(splitFile[6:], "/")

	// Log the message along with file name and line number
	fmt.Printf("%s:%d:::\n", fileName, line)
	log.Println(msg)
	fmt.Println()
}

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
	log.Printf(":::%s:%d:::\n%v\n", fileName, line, msg)
	fmt.Println()
}

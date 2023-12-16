package states

import (
	"log"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"github.com/joho/godotenv"
)

// User represents the metadata of an authenticated user
type User struct {
	ID       string
	Password []byte
	First    string
	Last     string
}

// Variables that holds all user and session details
var (
	MapUsers      = map[string]User{}
	MapSessions   = map[string]string{}
	MapAttendance = map[time.Time]map[string]time.Time{}
)

func init() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln("error initializing states::" + err.(error).Error())
		}
	}()

	logger.Println("Initializing envs...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	logger.Println("Success!")

	// init users from users.json
	logger.Println("Initializing users")
	// Can potentially panic if unable to read from file
	if err := db.Read("users.json", &MapUsers); err != nil {
		logger.Println(err)
	}
	logger.Println("Success!")

	// init attendance from attendance.json
	logger.Println("Initializing users")
	// Can potentially panic if unable to read from file
	if err := db.Read("attendance.json", &MapAttendance); err != nil {
		logger.Println(err)
	}
	logger.Println("Success!")
}

package states

import (
	"log"
	"sync"
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
	MapUsersMutex sync.Mutex
	MapUsers      = map[string]User{}

	MapSessionsMutex sync.Mutex
	MapSessions      = map[string]string{}

	MapAttendanceMutex sync.Mutex
	MapAttendance      = map[time.Time]map[string]time.Time{}
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

// GetMapUser is the thread-safe getter for values within MapUsers
func GetMapUser(userID string) (User, bool) {
	MapUsersMutex.Lock()
	defer MapUsersMutex.Unlock()
	user, ok := MapUsers[userID]
	return user, ok
}

// GetAllMapUsers is the thread-safe getter for MapUsers map
func GetAllMapUsers() map[string]User {
	MapUsersMutex.Lock()
	defer MapUsersMutex.Unlock()
	return MapUsers
}

// SetMapUser is the thread-safe setter for MapUsers
func SetMapUser(userID string, user User) {
	MapUsersMutex.Lock()
	defer MapUsersMutex.Unlock()
	MapUsers[userID] = user
}

// GetMapSession is the thread-safe getter for values within MapSessions
func GetMapSession(sessionID string) (string, bool) {
	MapSessionsMutex.Lock()
	defer MapSessionsMutex.Unlock()
	userID, ok := MapSessions[sessionID]
	return userID, ok
}

// GetAllMapSessions is the thread-safe getter for MapUsers map
func GetAllMapSessions() map[string]string {
	MapSessionsMutex.Lock()
	defer MapSessionsMutex.Unlock()
	return MapSessions
}

// SetMapSession is the thread-safe setter for MapSessions
func SetMapSession(sessionID, userID string) {
	MapSessionsMutex.Lock()
	defer MapSessionsMutex.Unlock()
	MapSessions[sessionID] = userID
}

// GetMapAttendanceOuter is the thread-safe getter for values within the outer MapAttendance map
func GetMapAttendanceOuter(dateTime time.Time) (map[string]time.Time, bool) {
	MapAttendanceMutex.Lock()
	defer MapAttendanceMutex.Unlock()

	innerMap, ok := MapAttendance[dateTime]
	if !ok {
		return nil, false
	}

	// Creating a copy to avoid concurrent map read and map write
	result := make(map[string]time.Time, len(innerMap))
	for k, v := range innerMap {
		result[k] = v
	}

	return result, true
}

// GetAllMapAttendanceOuter is the thread-safe getter for the outer MapAttendance map
func GetAllMapAttendanceOuter() map[time.Time]map[string]time.Time {
	MapAttendanceMutex.Lock()
	defer MapAttendanceMutex.Unlock()
	return MapAttendance
}

// GetMapAttendanceInner is the thread-safe getter for values within the inner MapAttendance map
func GetMapAttendanceInner(dateTime time.Time, userID string) (time.Time, bool) {
	MapAttendanceMutex.Lock()
	defer MapAttendanceMutex.Unlock()

	if userAttendance, ok := MapAttendance[dateTime]; ok {
		value, userExists := userAttendance[userID]
		return value, userExists
	}

	return time.Time{}, false
}

// SetMapAttendanceInner is the thread-safe setter for the inner MapAttendance map
func SetMapAttendanceInner(dateTime time.Time, userID string, value time.Time) {
	MapAttendanceMutex.Lock()
	defer MapAttendanceMutex.Unlock()

	if _, ok := MapAttendance[dateTime]; !ok {
		MapAttendance[dateTime] = make(map[string]time.Time)
	}

	MapAttendance[dateTime][userID] = value
}

// SetMapAttendanceOuter is the thread-safe setter for the outer MapAttendance map
func SetMapAttendanceOuter(dateTime time.Time, values map[string]time.Time) {
	MapAttendanceMutex.Lock()
	defer MapAttendanceMutex.Unlock()

	MapAttendance[dateTime] = values
}

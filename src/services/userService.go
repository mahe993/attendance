package services

import (
	"net/http"
	"sync"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/states"
	"attendance.com/src/templates"
	utils "attendance.com/src/util"
)

// UserPageVariables struct represents the variables used in the user page.
type UserPageVariables struct {
	User states.User
	Tab  string
}

// UserService struct provides methods for handling business logics for requests to the "/user" endpoint
type UserService struct {
	Variables   UserPageVariables
	VariablesMu sync.Mutex
}

// Usr is a global variable that provides access to the UserService methods
var (
	Usr UserService
)

// CheckIn handles the check-in process for a user.
// It guards if the user is already checked in, and if they are on the appropriate WIFI.
// It then updates the attendance.json file and redirects to the success page.
// If any error occurs during the check-in process, it recovers from the panic and redirects to the home page.
func (u *UserService) CheckIn(w http.ResponseWriter, r *http.Request) {
	// isCheckedIn potentially panics
	// Recover from panic and redirect to home page
	defer func() {
		if err := recover(); err != nil {
			logger.Println("error checking in::" + err.(error).Error())
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}()

	currUser := Auth.GetUser(r)
	if currUser.ID == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if user is already checked in
	if templates.IsCheckedIn(currUser.ID) != "" {
		http.Error(w, "You are already checked in", http.StatusForbidden)
		return
	}

	// Check if user is on appropriate WIFI
	ok, err := utils.ValidateClientIPHandler(r)
	if !ok || err != nil {
		logger.Println(err)
		http.Error(w, "Unable to check-in. You are not on the appropriate WIFI.", http.StatusForbidden)
		return
	}

	// Update attendance.json
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if _, ok := states.GetMapAttendanceOuter(today); !ok {
		states.SetMapAttendanceOuter(today, map[string]time.Time{})
	}
	states.SetMapAttendanceInner(today, currUser.ID, now)

	// Write MapAttendance state to database
	// Can potentially panic here if the database is not writable
	err = db.Write(states.GetAllMapAttendanceOuter(), "attendance.json")
	if err != nil {
		logger.Println(err)
	}

	http.Redirect(w, r, "/user/attendance/success", http.StatusFound)
}

// CheckInSuccess redirects the user to the home page with the "attendanceSuccess" form value set to "success".
// This is used to display a success message on the home page.
func (u *UserService) CheckInSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?attendanceSuccess=success", http.StatusFound)
}

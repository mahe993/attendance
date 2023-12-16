package services

import (
	"net/http"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/states"
	"attendance.com/src/templates"
	utils "attendance.com/src/util"
)

type UserPageVariables struct {
	User states.User
	Tab  string
}
type UserService struct {
	Variables UserPageVariables
}

var (
	Usr UserService
)

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
	ok, err := utils.ValidateIP()
	if !ok || err != nil {
		logger.Println(err)
		http.Error(w, "Unable to check-in. You are not on the appropriate WIFI.", http.StatusForbidden)
		return
	}

	// Update attendance.json
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if _, ok := states.MapAttendance[today]; !ok {
		states.MapAttendance[today] = map[string]time.Time{}
	}
	states.MapAttendance[today][currUser.ID] = now
	logger.Println("heREasdsad")
	// Write MapAttendance state to database
	// Can potentially panic here if the database is not writable
	err = db.Write(states.MapAttendance, "attendance.json")
	if err != nil {
		logger.Println(err)
	}

	http.Redirect(w, r, "/user/attendance/success", http.StatusFound)
}

func (u *UserService) CheckInSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?attendanceSuccess=success", http.StatusFound)
}

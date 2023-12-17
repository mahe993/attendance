/*
Package templates initializes and provides HTML templates for use in services.

The templates package includes a global variable Tpl, which is a pointer to a template.Template object. It also provides functions for template execution and rendering.

Initialization:

During package initialization, the Tpl variable is initialized with HTML templates and associated functions.

Note: The templates are assumed to be located in the "./templates/" directory and have a ".gohtml" extension.
*/
package templates

import (
	"html/template"
	"time"

	"attendance.com/src/logger"
	"attendance.com/src/states"
)

// AttendanceDetails struct represents details about a user's attendance, including check-in time and name
type AttendanceDetails struct {
	CheckInTime string
	Name        string
}

// CheckedInUsers is a map of date to map of user id to check in time
type CheckedInUsers map[string]map[string]AttendanceDetails

// Tpl is a pointer to a template.Template object that holds all initialized HTML templates
var Tpl *template.Template

func init() {
	logger.Println("Initializing templates...")
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"isCheckedIn": IsCheckedIn,
		"getCheckIns": GetCheckedInUsers,
	}).ParseGlob("./templates/*.gohtml"))
	logger.Println("Templates ready!")

}

// IsCheckedIn checks if a user is already checked in and returns the check-in time in a formatted string.
// If the user is not checked in, it returns an empty string.
func IsCheckedIn(id string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Check if user is already checked in
	if loggedInUsers, ok := states.GetMapAttendanceOuter(today); ok {
		if checkedInTime, ok := loggedInUsers[id]; ok {
			return checkedInTime.Format("Monday, 2 Jan 2006, 3:04:05 PM")
		}
		return ""
	}
	return ""
}

// GetCheckedInUsers retrieves a map of checked-in users within a specified date range.
// The date range is specified by the dateFrom and dateTo parameters, which are expected to be in the format "YYYY-MM-DD".
func GetCheckedInUsers(dateFrom string, dateTo string) CheckedInUsers {
	checkedInUsers := make(CheckedInUsers)
	if dateFrom == "" || dateTo == "" {
		return checkedInUsers
	}
	dateFromTime, err := time.ParseInLocation("2006-01-02", dateFrom, time.Now().Location())
	if err != nil {
		return checkedInUsers
	}
	dateToTime, err := time.ParseInLocation("2006-01-02", dateTo, time.Now().Location())
	if err != nil {
		return checkedInUsers
	}

	for dateFromTime.Before(dateToTime) || dateFromTime.Equal(dateToTime) {
		k := dateFromTime.Format("2006-01-02")
		if loggedInUsers, ok := states.GetMapAttendanceOuter(dateFromTime); ok {
			checkedInUsers[k] = make(map[string]AttendanceDetails)
			for id := range states.GetAllMapUsers() {
				if id == "admin" {
					continue
				}
				if checkedInTime, ok := loggedInUsers[id]; ok {
					if usr, ok := states.GetMapUser(id); ok {
						checkedInUsers[k][id] = AttendanceDetails{
							CheckInTime: checkedInTime.Format("Monday, 2 Jan 2006, 3:04:05 PM"),
							Name:        usr.First + " " + usr.Last,
						}
					}
				}
			}
		} else {
			checkedInUsers[k] = nil
		}
		dateFromTime = dateFromTime.AddDate(0, 0, 1)
	}

	return checkedInUsers
}

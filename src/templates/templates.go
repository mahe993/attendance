/*
Package templates initializes all the templates for use in services.
*/
package templates

import (
	"html/template"
	"time"

	"attendance.com/src/logger"
	"attendance.com/src/states"
)

type AttendanceDetails struct {
	CheckInTime string
	Name        string
}

// CheckedInUsers is a map of date to map of user id to check in time
type CheckedInUsers map[string]map[string]AttendanceDetails

var Tpl *template.Template

func init() {
	logger.Println("Initializing templates...")
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"isCheckedIn": IsCheckedIn,
		"getCheckIns": GetCheckedInUsers,
	}).ParseGlob("./templates/*.gohtml"))
	logger.Println("Templates ready!")

}

func IsCheckedIn(id string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Check if user is already checked in
	if loggedInUsers, ok := states.MapAttendance[today]; ok {
		if checkedInTime, ok := loggedInUsers[id]; ok {
			return checkedInTime.Format("Monday, 2 Jan 2006, 3:04:05 PM")
		}
		return ""
	}
	return ""
}

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
		if loggedInUsers, ok := states.MapAttendance[dateFromTime]; ok {
			checkedInUsers[k] = make(map[string]AttendanceDetails)
			for id := range states.MapUsers {
				if id == "admin" {
					continue
				}
				if checkedInTime, ok := loggedInUsers[id]; ok {
					checkedInUsers[k][id] = AttendanceDetails{
						CheckInTime: checkedInTime.Format("Monday, 2 Jan 2006, 3:04:05 PM"),
						Name:        states.MapUsers[id].First + " " + states.MapUsers[id].Last,
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

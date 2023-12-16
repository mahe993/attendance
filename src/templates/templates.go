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

var Tpl *template.Template

func init() {
	logger.Println("Initializing templates...")
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"isCheckedIn": IsCheckedIn,
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

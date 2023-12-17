package services

import (
	"log"
	"net/http"
	"sync"

	"attendance.com/src/states"
	"attendance.com/src/templates"
)

type MainPageVariables struct {
	User states.User
	Tab  string
}
type MainService struct {
	Variables   MainPageVariables
	VariablesMu sync.Mutex
}

var (
	MainPage MainService
)

func (p *MainService) Index(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)
	if currUser.ID == "admin" {
		http.Redirect(w, r, "/admin/overview", http.StatusFound)
		return
	}

	successTab := r.FormValue("attendanceSuccess")
	if successTab == "success" {
		if templates.IsCheckedIn(currUser.ID) == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	// Mutex lock to ensure thread-safe access to shared Variables field
	// Prevents race conditions when multiple requests concurrently read/write to Variables
	p.VariablesMu.Lock()
	p.Variables.User = currUser
	p.Variables.Tab = successTab

	err := templates.Tpl.ExecuteTemplate(w, "index", p.Variables)
	p.VariablesMu.Unlock()

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

package services

import (
	"log"
	"net/http"
	"sync"

	"attendance.com/src/states"
	"attendance.com/src/templates"
)

// MainPageVariables struct represents the variables used in the main page.
type MainPageVariables struct {
	User states.User
	Tab  string
}

// MainService struct provides methods for handling business logics for requests to the "/" endpoint
type MainService struct {
	Variables   MainPageVariables
	VariablesMu sync.Mutex
}

// MainPage is a global variable that provides access to the MainService methods
var (
	MainPage MainService
)

// Index handles the HTTP request for the main landing page.
// It redirects the user to the admin overview page if the current user is an admin.
// If the user is not an admin, it guards against users manually typing success routes if the "attendanceSuccess" form value is set to "success".
// It then updates the shared Variables field and executes the "index" template.
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

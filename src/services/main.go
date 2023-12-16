package services

import (
	"log"
	"net/http"

	"attendance.com/src/states"
	"attendance.com/src/templates"
)

type MainPageVariables struct {
	User states.User
	Tab  string
}
type MainService struct {
	Variables MainPageVariables
}

var (
	MainPage MainService
)

func (p *MainService) Index(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)
	p.Variables.User = currUser
	successTab := r.FormValue("attendanceSuccess")
	if successTab == "success" {
		if templates.IsCheckedIn(currUser.ID) == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	p.Variables.Tab = successTab

	err := templates.Tpl.ExecuteTemplate(w, "index", p.Variables)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

package services

import (
	"log"
	"net/http"

	"attendance.com/src/templates"
)

type MainPageVariables struct {
	User User
}
type MainService struct {
	Variables MainPageVariables
}

var (
	MainPage MainService
)

func (p *MainService) Index(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)
	if currUser != nil {
		p.Variables.User = *currUser
	} else {
		p.Variables.User = User{}
	}

	err := templates.Tpl.ExecuteTemplate(w, "index", p.Variables)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

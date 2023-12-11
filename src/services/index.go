package services

import (
	"fmt"
	"net/http"

	"attendance.com/templates"
)

func Index(w http.ResponseWriter, r *http.Request) {
	currUser := GetUser(w, r)
	err := templates.Tpl.ExecuteTemplate(w, "index.gohtml", currUser)
	fmt.Println(err)
}

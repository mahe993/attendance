package services

import (
	"log"
	"net/http"

	"attendance.com/src/templates"
)

type PageService struct{}

var (
	Page PageService
)

func (*PageService) Index(w http.ResponseWriter, r *http.Request) {
	currUser := GetUser(w, r)

	err := templates.Tpl.ExecuteTemplate(w, "index", currUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

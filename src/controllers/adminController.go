/*
Package controllers handles routing to specific services.
*/
package controllers

import (
	"net/http"
	"strings"

	"attendance.com/src/services"
)

type AdminController struct{}

var (
	Admin AdminController
)

func (*AdminController) Controller(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		Admin.POST(w, r)
	case http.MethodGet:
		Admin.GET(w, r)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		fallthrough
	default:
		http.NotFound(w, r)
	}
}

func (*AdminController) POST(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/admin")

	switch path {
	case "/upload":
		services.Admin.UploadStudentsList(w, r)
	case "/export":
		services.Admin.ExportAttendanceCSV(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (*AdminController) GET(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/admin")

	switch path {
	case "/upload":
		fallthrough
	case "/success":
		fallthrough
	case "/overview":
		services.Admin.Index(w, r)
	default:
		http.NotFound(w, r)
	}
}

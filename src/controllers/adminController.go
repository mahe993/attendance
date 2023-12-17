/*
Package controllers provides HTTP request handling for administrative operations.

The controllers package handles routing to specific services, facilitating the interaction with the attendance management system.

The controllers includes methods to handle HTTP methods such as POST and GET. Requests are routed to specific services based on the URL path.
*/
package controllers

import (
	"net/http"
	"strings"

	"attendance.com/src/services"
)

// AdminController handles HTTP request handling for administrative operations.
type AdminController struct{}

var (
	// Admin is an instance of AdminController.
	Admin AdminController
)

// Controller routes the HTTP request to the appropriate method based on the HTTP method.
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

// POST handles the HTTP POST request and routes it to the appropriate service based on the URL path.
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

// GET handles the HTTP GET request and routes it to the appropriate service based on the URL path.
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

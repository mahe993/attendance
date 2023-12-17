package controllers

import (
	"net/http"
	"strings"

	"attendance.com/src/services"
)

// UserController represents the controller for handling HTTP requests that are user related.
type UserController struct{}

var (
	// User is an instance of UserController.
	User UserController
)

// Controller routes the HTTP request to the appropriate method based on the HTTP method.
func (*UserController) Controller(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		User.POST(w, r)
	case http.MethodGet:
		User.GET(w, r)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		fallthrough
	default:
		http.NotFound(w, r)
	}
}

// POST routes HTTP POST requests to the appropriate services.
func (*UserController) POST(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/user")

	switch path {
	case "/attendance":
		services.Usr.CheckIn(w, r)
	default:
		http.NotFound(w, r)
	}
}

// GET routes HTTP GET requests to the appropriate services.
func (*UserController) GET(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/user")

	switch path {
	case "/attendance/success":
		services.Usr.CheckInSuccess(w, r)
	default:
		http.NotFound(w, r)
	}
}

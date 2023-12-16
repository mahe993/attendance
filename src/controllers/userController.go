/*
Package controllers handles routing to specific services.
*/
package controllers

import (
	"net/http"
	"strings"

	"attendance.com/src/services"
)

type UserController struct{}

var (
	User UserController
)

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

func (*UserController) POST(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/user")

	switch path {
	case "/attendance":
		services.Usr.CheckIn(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (*UserController) GET(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/user")

	switch path {
	case "/attendance/success":
		services.Usr.CheckInSuccess(w, r)
	default:
		http.NotFound(w, r)
	}
}

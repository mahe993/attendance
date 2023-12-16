/*
Package controllers handles routing to specific services.
*/
package controllers

import (
	"net/http"
	"strings"

	"attendance.com/src/services"
)

type AuthController struct{}

var (
	Auth AuthController
)

func (*AuthController) Controller(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		Auth.POST(w, r)
	case http.MethodGet:
		Auth.GET(w, r)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		fallthrough
	default:
		http.NotFound(w, r)
	}
}

func (*AuthController) POST(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/auth")

	switch path {
	case "/login":
		services.Auth.Login(w, r)
	case "/logout":
		services.Auth.Logout(w, r)
	case "/register":
		services.Auth.Register(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (*AuthController) GET(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/auth")

	switch path {
	case "/success":
		fallthrough
	case "/register":
		services.Auth.RegisterPage(w, r)
	default:
		http.NotFound(w, r)
	}
}

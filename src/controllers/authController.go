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
		fallthrough
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
	default:
		http.NotFound(w, r)
	}
}

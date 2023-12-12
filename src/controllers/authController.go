/*
Package controllers handles routing to specific services.
*/
package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"attendance.com/src/services"
)

func AuthController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(os.Stdout, "control:: "+r.URL.Path+"\n")

	path := strings.TrimPrefix(r.URL.Path, "/auth")

	fmt.Fprint(os.Stdout, "subpath:: "+path+"\n")

	switch path {
	case "":
		services.Login(w, r)
	case "/logout":
		services.Logout(w, r)
	default:
		http.NotFound(w, r)
	}

}

/*
Package router handles routing to specific controllers
*/
package router

import (
	"net/http"
	"strings"

	"attendance.com/src/controllers"
	"attendance.com/src/logger"
	"attendance.com/src/services"
)

func Routes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	logger.Println("route:: " + path + "\n")
	switch {
	case path == "/":
		services.MainPage.Index(w, r)
	case strings.HasPrefix(path, "/auth"):
		controllers.Auth.Controller(w, r)
	case strings.HasPrefix(path, "/admin"):
		if isAdmin := checkAdmin(w, r); !isAdmin {
			break
		}
		controllers.Admin.Controller(w, r)
	default:
		http.NotFound(w, r)
	}
}

func checkAdmin(w http.ResponseWriter, r *http.Request) bool {
	currUser := services.Auth.GetUser(r)

	if currUser.ID != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	return currUser.ID == "admin"
}

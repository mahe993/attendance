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
		if isAdmin := checkAuth(w, r, true); !isAdmin {
			break
		}
		controllers.Admin.Controller(w, r)
	case strings.HasPrefix(path, "/user"):
		if isAuthenticated := checkAuth(w, r, false); !isAuthenticated {
			break
		}
		controllers.User.Controller(w, r)
	// Handle static files
	case strings.HasSuffix(path, ".css"):
		http.ServeFile(w, r, "./templates/css/index.css")
	case strings.HasSuffix(path, ".js"):
		http.ServeFile(w, r, "./templates/scripts/script.js")
	default:
		http.NotFound(w, r)
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request, adminCheck bool) bool {
	currUser := services.Auth.GetUser(r)

	if currUser.ID == "" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if !adminCheck {
		return currUser.ID != ""
	}

	if currUser.ID != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	return currUser.ID == "admin"
}

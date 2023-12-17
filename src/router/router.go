/*
Package router handles HTTP request routing to specific controllers.

The router package provides a Routes function that routes the request to the appropriate controller based on the requested path.

The Routes function is responsible for routing HTTP requests to specific controllers based on the requested path. It includes logic to handle authentication, authorization, and serve static files.

Supported Paths:

The Routes function supports the following paths:

- /: Routes to the MainPage controller.
- /auth: Routes to the Auth controller.
- /admin: Routes to the Admin controller, performing admin authentication check.
- /user: Routes to the User controller, performing user authentication check.

Static Files:

Static files such as CSS and JavaScript are served for paths ending with ".css" and ".js" respectively.

Authentication and Authorization:

Routes can be protected by using the checkAuth function.
*/
package router

import (
	"net/http"
	"strings"

	"attendance.com/src/controllers"
	"attendance.com/src/logger"
	"attendance.com/src/services"
)

// The Routes function is responsible for routing HTTP requests to specific controllers based on the requested path.
// It includes logic to handle authentication, authorization, and serve static files.
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

// The checkAuth function is used to perform authentication and authorization checks based on the requested path.
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

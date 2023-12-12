/*
Package router handles routing to specific controllers
*/
package router

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"attendance.com/src/controllers"
	"attendance.com/src/services"
)

func Routes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Fprint(os.Stdout, "routes:: "+path+"\n")
	switch {
	case path == "/":
		services.Index(w, r)
	case strings.HasPrefix(path, "/auth"):
		controllers.Auth.Controller(w, r)
	default:
		http.NotFound(w, r)
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	currUser := services.GetUser(w, r)
	if currUser.ID == "" {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
	return currUser.ID == ""
}

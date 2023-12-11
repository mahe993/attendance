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
		controllers.AuthController(w, r)
	default:
		http.NotFound(w, r)
	}
}

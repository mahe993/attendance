/*
Package controllers handles routing to specific services.
*/
package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func AuthController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(os.Stdout, "control:: "+r.URL.Path+"\n")

	path := strings.TrimPrefix(r.URL.Path, "/auth")

	fmt.Fprint(os.Stdout, "subpath:: "+path+"\n")

	switch path {
	case "":
		fmt.Fprint(w, "working auth")
	case "/test":
		fmt.Fprint(w, "hehe")
	default:
		http.NotFound(w, r)
	}

}

/*
Package main is the entry point for the executable program.

Description:

	This program allows users to login and check-in to timestamp their attendance.
	Admins can upload a list of users through a .csv file.
	Users can login with their user ID which is made known to them by the admin, and is also recorded in the .csv.

Usage:

	$ go run main.go
*/
package main

import (
	"net/http"

	"attendance.com/src/router"
)

func main() {
	http.HandleFunc("/", router.Routes)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":5332", nil)
}

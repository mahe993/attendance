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
	"log"
	"net/http"

	"attendance.com/src/logger"
	"attendance.com/src/router"
)

func main() {
	http.HandleFunc("/", router.Routes)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./templates/scripts"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	logger.Println("Server listening at port :5332...")
	log.Fatal(http.ListenAndServe(":5332", nil))
	logger.Println("Server connection ended!")
}
